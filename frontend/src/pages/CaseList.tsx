import React, {useEffect, useMemo, useRef, useState} from 'react';
import {
    Box,
    Button,
    ButtonGroup,
    ChakraProvider,
    Flex,
    FormControl,
    FormLabel,
    Heading,
    HStack,
    Icon,
    IconButton,
    Input,
    Link,
    ListItem,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    OrderedList,
    Select,
    Table,
    Tbody,
    Td,
    Text,
    Textarea,
    Tr,
    UnorderedList,
    useDisclosure,
    useToast,
    VStack
} from '@chakra-ui/react';
import {DeleteIcon, DragHandleIcon} from '@chakra-ui/icons';
import {SlFolder} from "react-icons/sl";
import ReactMarkdown from 'react-markdown';
import gfm from 'remark-gfm';
import {Tree} from 'antd';
import {PiFilePlus, PiFolderSimplePlus} from "react-icons/pi";
import {useParams} from "react-router-dom";
import Header from "../components/Header";
import {useTranslation} from "react-i18next";
import {useApiRequest} from "../components/UseApiRequest";
import {DndProvider, useDrag, useDrop} from 'react-dnd';
import {HTML5Backend} from 'react-dnd-html5-backend';

interface TestCase {
    id: number;
    title: string;
    content: string;
    created_by: { id: number; name: string; };
    updated_by: { id: number; name: string; };
    milestone_id: number | null;
    milestone: Milestone | null;
    updated_by_id: number;
}

interface TestSuite {
    name: string;
    id: number;
    test_suites?: TestSuite[];
    test_cases?: TestCase[];
}

interface OnlyTestSuite {
    key: number;
    title: string
    children?: OnlyTestSuite[];
}

interface TestSuiteHeaderProps {
    title: string;
    testSuiteId: number;
    onAddCase: (testSuiteId: number) => void;
    onAddSuite: (testSuiteId: number | null) => void;
}

interface NewTestSuite {
    project_id: number;
    parent_id: number | null;
    name: string;
}

interface Milestone {
    id: number;
    title: string;
}

interface DraggableTestCaseProps {
    testCase: TestCase;
    index: number;
    testSuiteId: number;
}

interface DroppableTestSuiteProps {
    testSuite: TestSuite;
    onDrop: (testCase: TestCase, testSuite: TestSuite) => void;
    children: React.ReactNode;
}

interface DragItem {
    type: string;
    id: number;
    key: number;
    index: number;
    testCase: TestCase;
}

interface TreeDataItem {
    title: React.ReactNode;
    key: number;
    children?: TreeDataItem[];
}

interface DroppableTreeNodeProps {
    node: TreeDataItem;
    onDropTestCase: (item: DragItem, nodeKey: number) => void;
    onDropTestSuite: (item: DragItem, nodeKey: number) => void;
}

interface DroppableTreeProps {
    treeData: TreeDataItem[];
    onDropTestCase: (item: DragItem, nodeKey: number) => void;
    onDropTestSuite: (item: DragItem, nodeKey: number) => void;
}

interface TreeNodeProps {
    nodeData: TreeDataItem;
    onDropTestCase: (item: DragItem, nodeKey: number) => void;
    onDropTestSuite: (item: DragItem, nodeKey: number) => void
}

const CaseList: React.FC = () => {
    const {
        isOpen: isTestCaseAddModalOpen,
        onOpen: onTestCaseAddModalOpen,
        onClose: onTestCaseAddModalClose
    } = useDisclosure();
    const {
        isOpen: isTestSuiteAddModalOpen,
        onOpen: onTestSuiteAddModalOpen,
        onClose: onTestSuiteAddModalClose
    } = useDisclosure();
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const {isOpen: isDeleteModalOpen, onOpen: onDeleteModalOpen, onClose: onDeleteModalClose} = useDisclosure();
    const [deletingTestCaseId, setDeletingTestCaseId] = useState<number | null>(null);
    const [projectId, setProjectId] = useState<number>(0);
    const [editMode, setEditMode] = useState(false); // 編集モードの状態を追加
    const [testSuites, setTestSuites] = useState<TestSuite[]>([]);
    const [onlyTestSuites, setOnlyTestSuites] = useState<OnlyTestSuite[]>([]);
    const [selectedTestCase, setSelectedTestCase] = useState<TestCase | null>(null);
    const [newTestCase, setNewTestCase] = useState({
        project_id: projectId,
        test_suite_id: 0,
        title: '',
        content: '',
        created_by_id: user.id,
        updated_by_id: user.id
    });
    const [newTestSuite, setNewTestSuite] = useState<NewTestSuite>({
        project_id: projectId,
        parent_id: null,
        name: '',
    });
    const [milestones, setMilestone] = useState<Milestone[]>([]);
    const toast = useToast();
    const {project_code} = useParams();
    const {t} = useTranslation();
    const apiRequest = useApiRequest();
    const [hoverSuiteId, setHoverSuiteId] = useState<number | null>(null);
    const [hoverIndex, setHoverIndex] = useState<number | null>(null);
    const [hoverPosition, setHoverPosition] = useState<string | null>(null);

    const DraggableTestCase: React.FC<DraggableTestCaseProps> = ({testCase, index, testSuiteId}) => {
        const ref = useRef<HTMLDivElement>(null);

        const [{isDragging}, drag, preview] = useDrag(() => ({
            type: 'TEST_CASE',
            item: {type: 'TEST_CASE', id: testCase.id, index: index, testCase: testCase},
            collect: monitor => ({
                isDragging: monitor.isDragging(),
            }),
            // end: (item, monitor) => {
            //     const dropResult = monitor.getDropResult() as DragItem;
            //     if (item && dropResult) {
            //         const dragIndex = item.index;
            //         const hoverIndex = dropResult.index;
            //         if (dragIndex !== hoverIndex) {
            //             // onMove(dragIndex, hoverIndex);
            //         }
            //     }
            // },
        }));

        const [, drop] = useDrop({
            accept: 'TEST_CASE',
            hover: (item: DragItem, monitor) => {
                const hoverBoundingRect = ref.current?.getBoundingClientRect();
                const clientOffset = monitor.getClientOffset();

                if (!hoverBoundingRect || !clientOffset) {
                    return;
                }

                // ドロップ対象の中心のY座標を計算
                const hoverMiddleY = (hoverBoundingRect.bottom - hoverBoundingRect.top) / 2;
                // マウスポインタのY座標を取得
                const hoverClientY = clientOffset.y - hoverBoundingRect.top;

                // マウスポインタがアイテムの上半分にあるか下半分にあるかを判定
                if (hoverClientY < hoverMiddleY) {
                    setHoverPosition('upper')
                } else {
                    setHoverPosition('lower')
                }
                setHoverIndex(index);
                setHoverSuiteId(testSuiteId);
            },
        });

        drop(drag(preview(ref)));

        return (
            <Box ref={ref} p={1} my={1} bg={isDragging ? "gray.100" : ""} borderRadius="md"
                 opacity={isDragging ? 0.5 : 1}>
                <HStack>
                    <DragHandleIcon mr={2}/>
                    <Text>{testCase.title}</Text>
                </HStack>
            </Box>
        );
    };

    const DroppableTestSuite: React.FC<DroppableTestSuiteProps> = ({testSuite, onDrop, children}) => {
        const ref = useRef<HTMLDivElement>(null);
        const [, drop] = useDrop({
            accept: 'TEST_CASE',
            drop: (item: DragItem, monitor) => {
                if (monitor.isOver({shallow: true})) {
                    onDrop(item.testCase, testSuite);
                }
                return {...item, index: item.index}
            },
            collect: monitor => ({
                isOver: monitor.isOver({shallow: true}),
            }),
        });

        drop(ref)

        return (
            <div ref={ref}>
                {children}
            </div>
        );
    };

    const DroppableTreeNode: React.FC<DroppableTreeNodeProps> = ({node, onDropTestCase, onDropTestSuite}) => {
        const [, drop] = useDrop({
            accept: ['TEST_CASE', 'TEST_SUITE'],
            drop: (item: DragItem, monitor) => {
                if (item.type === 'TEST_CASE') {
                    onDropTestCase(item, node.key);
                } else if (item.type === 'TEST_SUITE') {
                    onDropTestSuite(item, node.key);
                }
            },
        });

        return (
            <div ref={drop}>
                {node.title}
                {node.children && node.children.map(childNode => (
                    <DroppableTreeNode key={childNode.key} node={childNode} onDropTestCase={onDropTestCase}
                                       onDropTestSuite={onDropTestSuite}/>
                ))}
            </div>
        );
    };

    const TreeNode: React.FC<TreeNodeProps> = ({nodeData, onDropTestCase, onDropTestSuite}) => {
        const [{isDragging}, drag] = useDrag({
            type: 'TEST_SUITE',
            item: {type: 'TEST_SUITE', key: nodeData.key},
            collect: monitor => ({
                isDragging: monitor.isDragging(),
            }),
        });

        const [{isOver}, drop] = useDrop({
            accept: ['TEST_CASE', 'TEST_SUITE'],
            drop: (item: DragItem, monitor) => {
                if (!monitor.didDrop()) {
                    if (item.type === 'TEST_CASE') {
                        onDropTestCase(item, nodeData.key);
                    } else if (item.type === 'TEST_SUITE') {
                        onDropTestSuite(item, nodeData.key);
                    }
                }
            },
            collect: monitor => ({
                isOver: monitor.isOver({shallow: true}),
            }),
        });

        const style = useMemo(() => ({
            backgroundColor: isOver ? '#BEE3F8' : isDragging ? '#F7FAFC' : 'white',
            // 他のスタイル設定
        }), [isOver, isDragging]);

        return (
            <div ref={(node) => drag(drop(node))} style={style}>
                {nodeData.title}
            </div>
        );
    };

    const DroppableTree: React.FC<DroppableTreeProps> = ({treeData, onDropTestCase, onDropTestSuite}) => {
        return (
            <Tree
                showLine
                defaultExpandAll
                treeData={treeData}
                titleRender={(nodeData) => <TreeNode nodeData={nodeData} onDropTestCase={onDropTestCase}
                                                     onDropTestSuite={onDropTestSuite}/>}
            />
        );
    }

    const onMoveTestCaseOnTestSuite = async (testCase: TestCase, testSuite: TestSuite) => {
        try {
            let testCases = testSuite.test_cases;

            if (!testCases || testCases.length == 0) {
                const response = await apiRequest(`/protected/cases/${testCase.id}`, {
                    method: 'PUT',
                    body: JSON.stringify({test_suite_id: testSuite.id}),
                });
                if (response.ok) {
                    setEditMode(false);
                    fetchTestCases();
                    fetchMilestones();
                    setHoverPosition(null);
                    setHoverSuiteId(null);
                    setHoverIndex(null);
                } else {
                    throw new Error(t('failed_to_update_test_case'));
                }
            } else {
                let movingTestCase = testCases.find(tc => tc.id === testCase.id);
                if (!movingTestCase) {
                    movingTestCase = testCase;
                    let newIndex = hoverIndex ? hoverIndex : 0;

                    // 新しい位置に挿入
                    if (newIndex >= testCases.length) {
                        testCases.push(movingTestCase);
                    } else {
                        if (hoverPosition == 'upper') {
                            testCases.splice(newIndex, 0, movingTestCase);
                        } else {
                            testCases.splice(newIndex - 1, 0, movingTestCase);
                        }
                    }
                } else {
                    // 同じスイート内での移動の場合
                    const filteredCases = testCases.filter(tc => tc.id !== testCase.id);
                    let newIndex = hoverIndex ? hoverIndex : 0;

                    if (newIndex >= filteredCases.length) {
                        filteredCases.push(movingTestCase);
                    } else {
                        filteredCases.splice(newIndex, 0, movingTestCase);
                    }
                    testCases = filteredCases;
                }

                console.log(JSON.stringify({test_suite_id: testSuite.id, test_cases: testCases.map((tc, index) => ({ test_case_id: tc.id, index }))}));
                const response = await apiRequest(`/protected/${project_code}/cases/bulk`, {
                    method: 'PUT',
                    body: JSON.stringify({test_suite_id: testSuite.id, test_cases: testCases.map((tc, index) => ({ test_case_id: tc.id, index }))}),
                });
                if (response.ok) {
                    setEditMode(false);
                    fetchTestCases();
                    fetchMilestones();
                    setHoverPosition(null);
                    setHoverSuiteId(null);
                    setHoverIndex(null);
                } else {
                    throw new Error(t('failed_to_update_test_case'));
                }
            }
        } catch (error) {
            let errorMessage = t('error_occurred');
            if (error instanceof Error) {
                errorMessage = error.message;
            }
            toast({
                title: t('error_occurred'),
                description: errorMessage,
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        }
    };

    const onMoveTestCase = async (testCaseId: number, targetTestSuiteId: number) => {
        try {
            const response = await apiRequest(`/protected/cases/${testCaseId}`, {
                method: 'PUT',
                body: JSON.stringify({test_suite_id: targetTestSuiteId}),
            });

            if (response.ok) {
                setEditMode(false);
                fetchTestCases();
                fetchMilestones();
            } else {
                throw new Error(t('failed_to_update_test_case'));
            }
        } catch (error) {
            let errorMessage = t('error_occurred');
            if (error instanceof Error) {
                errorMessage = error.message;
            }
            toast({
                title: t('error_occurred'),
                description: errorMessage,
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        }
    };

    const onMoveTestSuite = async (testSuiteId: number, targetTestSuiteId: number) => {
        try {
            const response = await apiRequest(`/protected/suites/${testSuiteId}`, {
                method: 'PUT',
                body: JSON.stringify({parent_id: targetTestSuiteId}),
            });

            if (response.ok) {
                setEditMode(false);
                fetchTestCases();
                fetchMilestones();
            } else {
                throw new Error(t('failed_to_update_test_case')); // TODO テストスイート用の文言に変える
            }
        } catch (error) {
            let errorMessage = t('error_occurred');
            if (error instanceof Error) {
                errorMessage = error.message;
            }
            toast({
                title: t('error_occurred'),
                description: errorMessage,
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        }
    };

    // APIからテストケースを取得
    const fetchTestCases = async () => {
        try {
            const response = await apiRequest(`/protected/${project_code}/cases`);
            const data = await response.json();
            setProjectId(data.project_id);
            setTestSuites(data.entities);
            setOnlyTestSuites(data.folders)
        } catch (error) {
            console.error('Error fetching TestCases:', error);
        }
    };

    const fetchMilestones = async () => {
        try {
            const response = await apiRequest(`/protected/${project_code}/milestones`);
            const data = await response.json();
            setMilestone(data.entities)
        } catch (error) {
            console.error('Error fetching Milestones:', error);
        }
    };

    useEffect(() => {
        fetchTestCases();
        fetchMilestones()
    }, [project_code]);

    const handleTestCaseClick = (testCase: TestCase) => {
        setSelectedTestCase(testCase);
    };

    const handleUpdateTestCase = async () => {
        if (selectedTestCase) {
            try {
                selectedTestCase.updated_by_id = user.id;
                const response = await apiRequest(`/protected/cases/${selectedTestCase.id}`, {
                    method: 'PUT',
                    body: JSON.stringify(selectedTestCase),
                });

                const updatedData = await response.json(); // 応答から最新のデータを取得
                setSelectedTestCase(updatedData);

                if (response.ok) {
                    toast({
                        title: t('test_case_updated'),
                        status: 'success',
                        duration: 5000,
                        isClosable: true,
                    });
                    setEditMode(false);
                    fetchTestCases();
                    fetchMilestones();
                } else {
                    throw new Error(t('failed_to_update_test_case'));
                }
            } catch (error) {
                let errorMessage = t('error_occurred');
                if (error instanceof Error) {
                    errorMessage = error.message;
                }
                toast({
                    title: t('error_occurred'),
                    description: errorMessage,
                    status: 'error',
                    duration: 5000,
                    isClosable: true,
                });
            }
        }
    };

    const handleAddTestCase = async () => {
        try {
            const updatedNewTestCase = {
                ...newTestCase,
                project_id: projectId
            };
            const response = await apiRequest(`/protected/cases`, {
                method: 'POST',
                body: JSON.stringify(updatedNewTestCase),
            });
            if (response.ok) {
                toast({
                    title: t('new_test_case_added'),
                    status: 'success',
                    duration: 5000,
                    isClosable: true,
                });
                onTestCaseAddModalClose();
                setNewTestCase({
                    project_id: projectId,
                    test_suite_id: 0,
                    title: '',
                    content: '',
                    created_by_id: user.id,
                    updated_by_id: user.id
                });
                fetchTestCases();
            } else {
                throw new Error(t('failed_to_add_new_test_case'));
            }
        } catch (error) {
            let errorMessage = t('error_occurred');
            if (error instanceof Error) {
                errorMessage = error.message;
            }
            toast({
                title: t('error_occurred'),
                description: errorMessage,
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        }
    };

    const handleAddTestSuite = async () => {
        try {
            const updatedNewTestSuite = {
                ...newTestSuite,
                project_id: projectId
            };
            const response = await apiRequest(`/protected/suites`, {
                method: 'POST',
                body: JSON.stringify(updatedNewTestSuite),
            });
            if (response.ok) {
                toast({
                    title: t('new_test_suite_added'),
                    status: 'success',
                    duration: 5000,
                    isClosable: true,
                });
                onTestSuiteAddModalClose();
                setNewTestCase({
                    project_id: projectId,
                    test_suite_id: 0,
                    title: '',
                    content: '',
                    created_by_id: user.id,
                    updated_by_id: user.id
                });
                fetchTestCases();
            } else {
                throw new Error(t('failed_to_add_new_test_suite'));
            }
        } catch (error) {
            let errorMessage = t('error_occurred');
            if (error instanceof Error) {
                errorMessage = error.message;
            }
            toast({
                title: t('error_occurred'),
                description: errorMessage,
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        }
    };

    const handleDeleteTestCase = async (id: number) => {
        try {
            // 削除確認が行われた後のAPIリクエスト
            const response = await apiRequest(`/protected/cases/${id}`, {
                method: 'DELETE'
            });
            // レスポンスがOKの場合、状態を更新してリストから削除
            if (response.ok) {
                // ユーザーに削除が成功したことを通知
                toast({
                    title: t('test_case_deleted'),
                    status: 'success',
                    duration: 5000,
                    isClosable: true,
                });
                fetchTestCases();
                fetchMilestones();
            } else {
                // レスポンスがOKでない場合、エラーをスロー
                throw new Error(t('failed_to_delete_test_case'));
            }
        } catch (error) {
            let errorMessage = t('error_occurred');
            if (error instanceof Error) {
                errorMessage = error.message;
            }
            toast({
                title: t('error_occurred'),
                description: errorMessage,
                status: 'error',
                duration: 5000,
                isClosable: true,
            });
        }
    };

    const onDelete = (testCaseId: React.SetStateAction<number | null>) => {
        setDeletingTestCaseId(testCaseId);
        onDeleteModalOpen();
    };

    const renderTestCases = (testCases: TestCase[], testSuiteId: number) => (
        <Table variant="simple">
            <Tbody>
                {testCases.map((testCase, index) => (
                    <React.Fragment key={testCase.id}>
                        {hoverIndex === index && hoverPosition === 'upper' && hoverSuiteId === testSuiteId && (
                            <Tr>
                                <Td colSpan={100} style={{height: '2px', backgroundColor: '#EDF2F7'}}/>
                            </Tr>
                        )}
                        <Tr id={testCase.id.toString()} cursor="pointer" _hover={{bg: "gray.100"}}
                            onClick={() => handleTestCaseClick(testCase)}>
                            <Td borderBottom="1px" borderColor="gray.200">
                                <DraggableTestCase testCase={testCase} index={index} testSuiteId={testSuiteId}/>
                            </Td>
                            <Td textAlign="right">
                                <IconButton
                                    aria-label="Delete test case"
                                    icon={<DeleteIcon/>}
                                    size="sm"
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        onDelete(testCase.id);
                                    }}
                                />
                            </Td>
                        </Tr>
                        {hoverIndex === index && hoverPosition === 'lower' && hoverSuiteId === testSuiteId && (
                            <Tr>
                                <Td colSpan={100} style={{height: '2px', backgroundColor: '#EDF2F7'}}/>
                            </Tr>
                        )}
                    </React.Fragment>
                ))}
            </Tbody>
        </Table>
    );

    const TestSuiteHeader: React.FC<TestSuiteHeaderProps> = ({title, testSuiteId, onAddCase, onAddSuite}) => {
        return (
            <Flex justifyContent="space-between" alignItems="center" mb={4}>
                <Flex alignItems="center">
                    <Icon as={SlFolder} mr={2}/>
                    <Text fontSize="lg" fontWeight="bold">{title}</Text>
                </Flex>
                <Flex>
                    <IconButton
                        aria-label={t('add_case')}
                        icon={<PiFilePlus/>}
                        colorScheme="gray"
                        size="sm"
                        onClick={() => onAddCase(testSuiteId)}
                    />
                    {/* New Button for Adding Folder */}
                    <IconButton
                        aria-label={t('add_test_suite')}
                        icon={<PiFolderSimplePlus/>}
                        colorScheme="gray"
                        size="sm"
                        ml={2}
                        onClick={() => onAddSuite(testSuiteId)}
                    />
                </Flex>
            </Flex>
        );
    };

    const handleTestCaseDropOnTree = (item: DragItem, targetNodeKey: number) => {
        const testCaseId = item.testCase.id;
        onMoveTestCase(testCaseId, targetNodeKey);
    };

    const handleTestSuiteDropOnTree = (item: DragItem, targetNodeKey: number) => {
        const testSuiteId = item.key;
        onMoveTestSuite(testSuiteId, targetNodeKey);
    };

    const renderTestSuites = (suites: TestSuite[], level = 0) => (
        suites.map(suite => (
            <DroppableTestSuite key={suite.id} testSuite={suite} onDrop={onMoveTestCaseOnTestSuite}>
                <Box mb={4} pl={`${level * 2}em`}>
                    <TestSuiteHeader
                        title={suite.name}
                        testSuiteId={suite.id}
                        onAddCase={(id) => {
                            setNewTestCase({...newTestCase, test_suite_id: id});
                            onTestCaseAddModalOpen();
                        }}
                        onAddSuite={(id) => {
                            setNewTestSuite({...newTestSuite, parent_id: id});
                            onTestSuiteAddModalOpen();
                        }}
                    />
                    <Box mb={4}>
                        {suite.test_cases && suite.test_cases.length > 0 && renderTestCases(suite.test_cases, suite.id)}
                    </Box>
                    {/* 再帰的に子のTestSuiteをレンダリング */}
                    {suite.test_suites && suite.test_suites.length > 0 && (
                        <Box>
                            {renderTestSuites(suite.test_suites, level + 1)}
                        </Box>
                    )}
                </Box>
            </DroppableTestSuite>
        ))
    );

    return (
        <DndProvider backend={HTML5Backend}>
            <ChakraProvider>
                <Header project_code={project_code} is_show_menu={true}/>
                <Flex h="100vh">
                    <Box w="20%" p={5} borderRight="1px" borderColor="gray.200" pt="6rem">
                        <DroppableTree treeData={onlyTestSuites} onDropTestCase={handleTestCaseDropOnTree}
                                       onDropTestSuite={handleTestSuiteDropOnTree}/>
                    </Box>
                    <Box w="50%" p={5} overflowY="auto" borderRight="1px" borderColor="gray.200" pt="6rem">
                        <Box p={5}>
                            <IconButton
                                aria-label={t('add_test_suite')}
                                icon={<PiFolderSimplePlus/>}
                                colorScheme="gray"
                                size="sm"
                                ml={2}
                                onClick={() => {
                                    setNewTestSuite({
                                        project_id: projectId,
                                        parent_id: null,
                                        name: ''
                                    });
                                    onTestSuiteAddModalOpen(); // モーダルを開く
                                }}
                                mb={2}
                            />
                            {renderTestSuites(testSuites)}
                            <Modal isOpen={isTestCaseAddModalOpen} onClose={onTestCaseAddModalClose}>
                                <ModalOverlay/>
                                <ModalContent>
                                    <ModalCloseButton/>
                                    <ModalHeader>{t('add_new_test_case')}</ModalHeader>
                                    <ModalBody>
                                        <VStack spacing={4}>
                                            <FormControl>
                                                <FormLabel>{t('title')}</FormLabel>
                                                <Input value={newTestCase.title}
                                                       onChange={(e) => setNewTestCase({
                                                           ...newTestCase,
                                                           title: e.target.value
                                                       })}/>
                                            </FormControl>
                                            <FormControl>
                                                <FormLabel>{t('content')}</FormLabel>
                                                <Textarea value={newTestCase.content}
                                                          height={200}
                                                          onChange={(e) => setNewTestCase({
                                                              ...newTestCase,
                                                              content: e.target.value
                                                          })}/>
                                            </FormControl>
                                        </VStack>
                                    </ModalBody>
                                    <ModalFooter>
                                        <Button variant="outline" mr={3} onClick={onTestCaseAddModalClose}>
                                            {t('cancel')}
                                        </Button>
                                        <Button colorScheme="blue" onClick={handleAddTestCase}>
                                            {t('add')}
                                        </Button>
                                    </ModalFooter>
                                </ModalContent>
                            </Modal>
                            <Modal isOpen={isTestSuiteAddModalOpen} onClose={onTestSuiteAddModalClose}>
                                <ModalOverlay/>
                                <ModalContent>
                                    <ModalHeader>{t('add_new_test_suite')}</ModalHeader>
                                    <ModalCloseButton/>
                                    <ModalBody>
                                        <FormControl>
                                            <FormLabel>{t('name')}</FormLabel>
                                            <Input value={newTestSuite.name}
                                                   onChange={(e) => setNewTestSuite({
                                                       ...newTestSuite,
                                                       name: e.target.value
                                                   })}/>
                                        </FormControl>
                                    </ModalBody>
                                    <ModalFooter>
                                        <Button colorScheme="blue" mr={3} onClick={handleAddTestSuite}>
                                            {t('add')}
                                        </Button>
                                        <Button onClick={onTestSuiteAddModalClose}>{t('cancel')}</Button>
                                    </ModalFooter>
                                </ModalContent>
                            </Modal>
                            <Modal isOpen={isDeleteModalOpen} onClose={onDeleteModalClose}>
                                <ModalOverlay/>
                                <ModalContent>
                                    <ModalHeader>{t('delete_test_case')}</ModalHeader>
                                    <ModalCloseButton/>
                                    <ModalBody>
                                        <Text>{t('confirm_delete')}</Text>
                                    </ModalBody>
                                    <ModalFooter>
                                        <Button colorScheme="red" onClick={() => {
                                            if (deletingTestCaseId) handleDeleteTestCase(deletingTestCaseId);
                                            onDeleteModalClose();
                                        }}>
                                            {t('delete')}
                                        </Button>
                                        <Button ml={2} onClick={onDeleteModalClose}>{t('cancel')}</Button>
                                    </ModalFooter>
                                </ModalContent>
                            </Modal>
                        </Box>
                    </Box>
                    <Box w="30%" p={5} pt="6rem">
                        {selectedTestCase && (
                            <VStack align="start">
                                {editMode ? (
                                    // 編集フォームを表示
                                    <>
                                        <FormControl>
                                            <FormLabel>{t('title')}</FormLabel>
                                            <Input
                                                value={selectedTestCase?.title}
                                                onChange={(e) => setSelectedTestCase({
                                                    ...selectedTestCase,
                                                    title: e.target.value
                                                })}
                                            />
                                        </FormControl>
                                        <FormControl>
                                            <FormLabel>{t('content')}</FormLabel>
                                            <Textarea
                                                height={200}
                                                value={selectedTestCase?.content}
                                                onChange={(e) => setSelectedTestCase({
                                                    ...selectedTestCase,
                                                    content: e.target.value
                                                })}
                                            />
                                        </FormControl>
                                        <FormControl>
                                            <FormLabel>{t('milestone')}</FormLabel>
                                            <Select
                                                placeholder={t('select_milestone')}
                                                onChange={(e) => setSelectedTestCase({
                                                    ...selectedTestCase,
                                                    milestone_id: parseInt(e.target.value, 10)
                                                })}
                                                value={selectedTestCase.milestone_id !== null ? selectedTestCase.milestone_id : ''}
                                            >
                                                {milestones.map((milestone) => (
                                                    <option key={milestone.id} value={milestone.id}>
                                                        {milestone.title}
                                                    </option>
                                                ))}
                                            </Select>
                                        </FormControl>
                                        <ButtonGroup size="sm">
                                            <Button colorScheme="blue"
                                                    onClick={handleUpdateTestCase}>{t('update')}</Button>
                                            <Button onClick={() => setEditMode(false)}>{t('cancel')}</Button>
                                        </ButtonGroup>
                                    </>
                                ) : (
                                    <>
                                        <VStack align="start">
                                            <Flex justify="space-between">
                                                <Heading as="h3" size="md">{selectedTestCase.title}</Heading>
                                                <Button size="sm" onClick={() => {
                                                    if (selectedTestCase && selectedTestCase.milestone) {
                                                        setSelectedTestCase({
                                                            ...selectedTestCase,
                                                            milestone_id: selectedTestCase.milestone.id
                                                        });
                                                    }
                                                    setEditMode(true);
                                                }} ml={2}>{t('edit')}</Button>
                                            </Flex>
                                        </VStack>
                                        <ReactMarkdown remarkPlugins={[gfm]} components={{
                                            h1: ({node, ...props}) => <Heading as="h1" size="xl" mt={6}
                                                                               mb={4} {...props} />,
                                            h2: ({node, ...props}) => <Heading as="h2" size="lg" mt={5}
                                                                               mb={3} {...props} />,
                                            h3: ({node, ...props}) => <Heading as="h3" size="md" mt={4}
                                                                               mb={2} {...props} />,
                                            h4: ({node, ...props}) => <Heading as="h4" size="sm" mt={3}
                                                                               mb={1} {...props} />,
                                            p: ({node, ...props}) => <Text mt={2} mb={2} {...props} />,
                                            a: ({node, ...props}) => <Link color="teal.500" isExternal {...props} />,
                                            ul: ({node, ...props}) => <UnorderedList mt={2} mb={2} {...props} />,
                                            ol: ({node, ...props}) => <OrderedList mt={2} mb={2} {...props} />,
                                            li: ({node, ...props}) => <ListItem {...props} />,
                                            em: ({node, ...props}) => <Text as="em" {...props} />,
                                            strong: ({node, ...props}) => <Text as="strong" {...props} />,
                                        }}>
                                            {selectedTestCase?.content}
                                        </ReactMarkdown>
                                        <Text fontSize="md" color="gray.600">
                                            {`Milestone: ${
                                                milestones.find(milestone => milestone.id === selectedTestCase.milestone?.id)?.title || 'None'
                                            }`}
                                        </Text>
                                    </>
                                )}
                            </VStack>
                        )}
                    </Box>
                </Flex>
            </ChakraProvider>
        </DndProvider>
    );
};

export default CaseList;
