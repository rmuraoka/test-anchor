import React, {useEffect, useState} from 'react';
import {
    Box,
    Button,
    ButtonGroup,
    ChakraProvider,
    Flex,
    FormControl,
    FormLabel,
    Heading,
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
    OrderedList, Select,
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
import {DeleteIcon} from '@chakra-ui/icons';
import {SlFolder} from "react-icons/sl";
import ReactMarkdown from 'react-markdown';
import gfm from 'remark-gfm';
import {Tree} from 'antd';
import {PiFilePlus, PiFolderSimplePlus} from "react-icons/pi";
import {useParams} from "react-router-dom";
import Header from "../components/Header";
import {useTranslation} from "react-i18next";
import {useApiRequest} from "../components/UseApiRequest";

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
    const { t } = useTranslation();
    const apiRequest = useApiRequest();

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
                setTestSuites(prevSuites => {
                    // 削除するテストケースを除外した新しいテストスイートのリストを作成
                    return prevSuites.map(suite => ({
                        ...suite,
                        test_cases: suite.test_cases?.filter(testCase => testCase.id !== id),
                    }));
                });

                // ユーザーに削除が成功したことを通知
                toast({
                    title: t('test_case_deleted'),
                    status: 'success',
                    duration: 5000,
                    isClosable: true,
                });
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

    const renderTestCases = (testCases: TestCase[]) => (
        <Table variant="simple">
            <Tbody>
                {testCases.map(testCase => (
                    <Tr cursor="pointer" _hover={{bg: "gray.100"}} onClick={() => handleTestCaseClick(testCase)}>
                        <Td borderBottom="1px" borderColor="gray.200">{testCase.title}</Td>
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

    const renderTestSuites = (suites: TestSuite[], level = 0) => (
        suites.map(suite => (
            <Box key={suite.name} mb={4} pl={`${level * 2}em`}>
                <TestSuiteHeader
                    title={suite.name}
                    testSuiteId={suite.id} // TestSuiteのIDを渡す
                    onAddCase={(id) => {
                        setNewTestCase({...newTestCase, test_suite_id: id}); // テストスイートのIDを新しいテストケースのステートに設定
                        onTestCaseAddModalOpen();
                    }}
                    onAddSuite={(id) => {
                        setNewTestSuite({...newTestSuite, parent_id: id}); // テストスイートのIDを新しいテストケースのステートに設定
                        onTestSuiteAddModalOpen();
                    }}
                />
                <Box mb={4}>
                    {suite.test_cases && suite.test_cases.length > 0 && renderTestCases(suite.test_cases)}
                </Box>
                {suite.test_suites && suite.test_suites.length > 0 && renderTestSuites(suite.test_suites, level + 1)}
            </Box>
        ))
    );

    return (
        <ChakraProvider>
            <Header project_code={project_code} is_show_menu={true}/>
            <Flex h="100vh">
                <Box w="20%" p={5} borderRight="1px" borderColor="gray.200" pt="6rem">
                    <Tree
                        showLine
                        defaultExpandAll
                        treeData={onlyTestSuites}
                    />
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
                                        <Button colorScheme="blue" onClick={handleUpdateTestCase}>{t('update')}</Button>
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
    );
};

export default CaseList;
