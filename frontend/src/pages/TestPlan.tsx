import React, {useEffect, useState} from 'react';
import {
    Box,
    Button,
    ButtonGroup,
    ChakraProvider,
    Checkbox,
    Container,
    Editable,
    EditableInput,
    EditablePreview,
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
    Spinner,
    Table,
    Tbody,
    Td,
    Text,
    Th,
    Thead,
    Tr,
    UnorderedList,
    useDisclosure,
    useEditableControls,
    useToast,
    useToken,
    VStack
} from '@chakra-ui/react';
import {useParams} from "react-router-dom";
import {
    ArcElement,
    CategoryScale,
    Chart as ChartJS,
    ChartOptions,
    Legend,
    LinearScale,
    PieController,
    Title,
    Tooltip
} from 'chart.js';
import {Pie} from "react-chartjs-2";
import {SlFolder} from "react-icons/sl";
import {Tree} from "antd";
import ReactMarkdown from "react-markdown";
import gfm from "remark-gfm";
import Header from "../components/Header";
import {useTranslation} from "react-i18next";
import {useApiRequest} from "../components/UseApiRequest";
import {CheckIcon, CloseIcon, EditIcon, TimeIcon} from "@chakra-ui/icons";

ChartJS.register(
    ArcElement,
    Tooltip,
    Legend,
    CategoryScale,
    LinearScale,
    Title,
    PieController
);

interface TestPlan {
    id: number;
    title: string;
    status: string;
    started_at: string | null;
    completed_at: string | null;
    test_runs: TestRun[];
    charts: Chart[];
    created_by: User;
    updated_by: User;
}

interface TestRun {
    id: number;
    title: string;
    started_at: string | null;
    completed_at: string | null;
    test_case_ids: number[];
    count: number;
    created_by: User;
    updated_by: User;
}

interface Chart {
    "name": string;
    "color": string;
    "count": number;
}

interface User {
    id: number;
    name: string;
}

interface TestCase {
    id: number;
    title: string;
    content: string;
    created_by: { id: number; name: string; };
    updated_by: { id: number; name: string; };
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
    testSuite: TestSuite;
    selectedTestCaseIds: number[];
    setSelectedTestCaseIds: React.Dispatch<React.SetStateAction<number[]>>;
}

interface TestPlanUpdate {
    status: string;
    started_at?: string;
    completed_at?: string | null;
}

const TestPlan: React.FC = () => {
    const [projectId, setProjectId] = useState<number>(0);
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const [testRuns, setTestRun] = useState<TestRun[]>([]);
    const [charts, setCharts] = useState<Chart[]>([]);
    const {isOpen, onOpen, onClose} = useDisclosure();
    const [isLoading, setIsLoading] = useState(true);
    const [testSuites, setTestSuites] = useState<TestSuite[]>([]);
    const [onlyTestSuites, setOnlyTestSuites] = useState<OnlyTestSuite[]>([]);
    const [selectedTestCase, setSelectedTestCase] = useState<TestCase | null>(null);
    const [newTestRun, setNewTestRun] = useState({
        project_id: 1,
        test_plan_id: 1,
        title: '',
        status: 'NotExecuted',
        started_at: null,
        completed_at: null,
        created_by_id: user.id,
        updated_by_id: user.id
    });
    const [searchTerm, setSearchTerm] = useState('');
    const toast = useToast();
    const {project_code, test_plan_id} = useParams();
    const [selectedTestRun, setSelectedTestRun] = useState<TestRun | null>(null); // 編集するテストラン
    const {isOpen: isEditModalOpen, onOpen: onEditModalOpen, onClose: onEditModalClose} = useDisclosure();
    const [selectedTestCaseIds, setSelectedTestCaseIds] = useState<number[]>([]); // 選択されたテストケースのID
    const {t} = useTranslation();
    const apiRequest = useApiRequest();
    const [selectedCount, setSelectedCount] = useState(0);
    const isOverLimit = selectedCount > 10000;
    const [testPlanStatus, setTestPlanStatus] = useState<string>('');

    // APIからテストケースを取得
    const fetchTestRuns = async () => {
        try {
            const response = await apiRequest(`/protected/plans/${test_plan_id}`);
            const data = await response.json();
            setProjectId(data.project_id)
            setTestRun(data.test_runs)
            setTestPlanStatus(data.status);
            setCharts(data.charts)
            setIsLoading(false);
        } catch (error) {
            console.error('Error fetching TestPlans:', error);
        }
    };

    const fetchTestCases = async () => {
        try {
            const response = await apiRequest(`/protected/${project_code}/cases`);
            const data = await response.json();
            setProjectId(data.project_id)
            setTestSuites(data.entities);
            setOnlyTestSuites(data.folders)
        } catch (error) {
            console.error('Error fetching TestCases:', error);
        }
    };

    useEffect(() => {
        fetchTestRuns();
        fetchTestCases();
    }, []);

    const handleCheckboxChange = (testCaseId: number, isChecked: boolean) => {
        setSelectedTestCaseIds(prev => {
            if (isChecked) {
                setSelectedCount(selectedCount + 1);
                // チェックされた場合、IDを追加
                return [...prev, testCaseId];
            } else {
                setSelectedCount(selectedCount - 1);
                // チェックが外れた場合、IDを削除
                return prev.filter(id => id !== testCaseId);
            }
        });
    };

    const handleUpdateButtonClick = async () => {
        try {
            const response = await apiRequest(`/protected/runs/cases/bulk`, {
                method: 'POST',
                body: JSON.stringify({
                    test_case_ids: selectedTestCaseIds,
                    test_run_id: selectedTestRun?.id
                })
            });

            if (!response.ok) {
                throw new Error('APIリクエストに失敗しました');
            }

            // 成功メッセージを表示
            toast({
                title: t('test_run_updated'),
                status: 'success',
                duration: 5000,
                isClosable: true,
            });
            fetchTestRuns();
            fetchTestCases();
            onEditModalClose();
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

    const handleTestPlanStatusChange = async (newStatus: string, currentStatus: string) => {
        try {
            const now = new Date().toISOString();
            let updateData: TestPlanUpdate = {
                status: newStatus
            };

            if (newStatus === 'InProcess' && currentStatus !== 'Completed') {
                updateData.started_at = now;
            } else if (newStatus == 'InProcess' && currentStatus == 'Completed') {
                updateData.completed_at = null;
            } else if (newStatus === 'Completed') {
                updateData.completed_at = now;
            }

            const response = await apiRequest(`/protected/plans/${test_plan_id}`, {
                method: 'PUT',
                body: JSON.stringify(updateData)
            });

            if (!response.ok) {
                throw new Error(t('status_update_failed'));
            }
            setTestPlanStatus(newStatus);
            fetchTestRuns();
            fetchTestCases();
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

    const handleTestCaseClick = (testCase: TestCase) => {
        setSelectedTestCase(testCase);
    };

    const handleAddTestRun = async () => {
        try {
            const updatedNewTestRun = {
                ...newTestRun,
                project_id: projectId,
                test_plan_id: Number(test_plan_id)
            };
            const response = await apiRequest(`/protected/runs`, {
                method: 'POST',
                body: JSON.stringify(updatedNewTestRun)
            });
            if (response.ok) {
                toast({
                    title: t('new_test_run_added'),
                    status: 'success',
                    duration: 5000,
                    isClosable: true,
                });
                setNewTestRun({
                    project_id: projectId,
                    test_plan_id: Number(test_plan_id),
                    title: '',
                    status: 'NotExecuted',
                    started_at: null,
                    completed_at: null,
                    created_by_id: user.id,
                    updated_by_id: user.id
                });
                onClose();
                fetchTestRuns();
            } else {
                throw new Error(t('failed_to_add_new_test_run'));
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
    const handleSaveTitle = async () => {
        try {
            if (!selectedTestRun) {
                return;
            }

            const response = await apiRequest(`/protected/runs/${selectedTestRun.id}`, {
                method: 'PUT',
                body: JSON.stringify({title: selectedTestRun.title, updated_by_id: user.id})
            });
            if (response.ok) {
                toast({
                    title: t('test_run_updated'),
                    status: 'success',
                    duration: 5000,
                    isClosable: true,
                });
                fetchTestRuns();
            } else {
                throw new Error(t('failed_to_update_test_run'));
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

    const handleEditTestRun = (testRun: TestRun) => {
        setSelectedTestCaseIds(testRun.test_case_ids)
        setSelectedCount(testRun.test_case_ids.length)
        setSelectedTestRun(testRun);
        onEditModalOpen();
    };

    const renderTestCases = (testCases: TestCase[]) => (
        <Table variant="simple">
            <Tbody>
                {testCases.map(testCase => (
                    <Tr cursor="pointer" _hover={{bg: "gray.100"}} onClick={() => handleTestCaseClick(testCase)}>
                        <Td borderBottom="1px" borderColor="gray.200" width="30px" paddingY="0">
                            <Checkbox
                                onChange={e => handleCheckboxChange(testCase.id, e.target.checked)}
                                isChecked={selectedTestCaseIds.includes(testCase.id)}
                            />
                        </Td>
                        <Td borderBottom="1px" borderColor="gray.200" paddingY="0.6em"><Text fontSize="sm">{testCase.title}</Text></Td>
                    </Tr>
                ))}
            </Tbody>
        </Table>
    );

    const MenuBar = () => {
        return (
            <Flex
                as="nav"
                position="fixed"
                bottom="0"
                left="0"
                right="0"
                borderTopWidth="1px"
                padding="1rem"
                justifyContent="start"
                bg="white"
                zIndex="sticky"
            >
                {renderStatusButton()}
            </Flex>
        );
    };

    const renderStatusButton = () => {
        switch (testPlanStatus) {
            case 'NotExecuted':
                return <Button leftIcon={<TimeIcon/>}
                               onClick={() => handleTestPlanStatusChange('InProcess', testPlanStatus)}>{t('start_test')}</Button>;
            case 'InProcess':
                return <Button leftIcon={<TimeIcon/>}
                               onClick={() => handleTestPlanStatusChange('Completed', testPlanStatus)}>{t('complete')}</Button>;
            case 'Completed':
                return <Button leftIcon={<TimeIcon/>}
                               onClick={() => handleTestPlanStatusChange('InProcess', testPlanStatus)}>{t('return_to_in_process')}</Button>;
            default:
                return null;
        }
    };

    const TestSuiteHeader: React.FC<TestSuiteHeaderProps> = ({testSuite, selectedTestCaseIds, setSelectedTestCaseIds}) => {
        const isAllChecked = testSuite.test_cases?.every(testCase => selectedTestCaseIds.includes(testCase.id));
        const handleTestSuiteCheckboxChange = (checked: boolean) => {
            if (checked) {
                const newSelectedIds = Array.from(new Set([...selectedTestCaseIds, ...testSuite.test_cases!.map(tc => tc.id)]));
                const count = testSuite.test_cases ? testSuite.test_cases.length : 0;
                setSelectedCount(selectedCount + count);
                setSelectedTestCaseIds(newSelectedIds);
            } else {
                const newSelectedIds = selectedTestCaseIds.filter(id => !testSuite.test_cases!.some(tc => tc.id === id));
                const count = testSuite.test_cases ? testSuite.test_cases.length : 0;
                setSelectedCount(selectedCount - count);
                setSelectedTestCaseIds(newSelectedIds);
            }
        };

        return (
            <Flex justifyContent="space-between" alignItems="center" mb={1}>
                <Flex
                    alignItems="center"
                    id={'testSuite'+testSuite.id.toString()}
                >
                    <Checkbox
                        isChecked={isAllChecked}
                        onChange={e => handleTestSuiteCheckboxChange(e.target.checked)}
                        mr={2}
                    />
                    <Text fontSize="lg" fontWeight="bold">{testSuite.name}</Text>
                </Flex>
            </Flex>
        );
    };

    const renderTestSuites = (suites: TestSuite[]) => (
        suites.map(suite => (
            <Box key={suite.name} mb={4} pl={`2em`}>
                <TestSuiteHeader testSuite={suite} selectedTestCaseIds={selectedTestCaseIds} setSelectedTestCaseIds={setSelectedTestCaseIds}/>
                <Box mb={4}>
                    {suite.test_cases && suite.test_cases.length > 0 && renderTestCases(suite.test_cases)}
                </Box>
                {suite.test_suites && suite.test_suites.length > 0 && renderTestSuites(suite.test_suites)}
            </Box>
        ))
    );

    function Chart() {
        const [red300, green300, gray300, orange300, yellow300, teal300, blue300, cyan300, purple300, pink300] =
            useToken("colors", ["red.300", "green.300", "gray.300", "orange.300", "yellow.300", "teal.300", "blue.300", "cyan.300", "purple.300", "pink.300"]);
        const hasData = charts.some(chart => chart.count > 0);
        if (!hasData) {
            return null; // データがない場合は何も表示しない
        }

        const data = {
            labels: charts.map(chart => chart.name),
            datasets: [
                {
                    data: charts.map(chart => chart.count),
                    backgroundColor: charts.map(chart => {
                        switch (chart.color) {
                            case "red":
                                return red300;
                            case "orange":
                                return orange300;
                            case "yellow":
                                return yellow300;
                            case "green":
                                return green300;
                            case "teal":
                                return teal300;
                            case "blue":
                                return blue300;
                            case "cyan":
                                return cyan300;
                            case "purple":
                                return purple300;
                            case "pink":
                                return pink300;
                            case "gray":
                            default:
                                return gray300; // デフォルトの色
                        }
                    }),
                    borderWidth: 0,
                },
            ],
        };

        const options: ChartOptions<"pie"> = {
            responsive: true,
            animation: {
                duration: 0
            },
            plugins: {
                legend: {
                    display: false,
                },
                title: {
                    display: false,
                    text: '',
                },
            },
        };

        return (
            <Container maxW="container.xl" py={10}>
                <Box border="1px" borderColor="gray.200" borderRadius="md" p={4}>
                    <Flex justifyContent="center" alignItems="center" mb={4}>
                        <Box width={"300px"} mr={3}>
                            <Pie data={data} options={options}/>
                        </Box>
                        <Box width={"10rem"}/>
                        <Flex direction="column" alignItems="flex-start">
                            {charts.map((chart, index) => (
                                <Box key={index} bg={`${chart.color}.100`} p={3} borderRadius="md"
                                     mt={index > 0 ? 2 : 0}
                                     width="300px">
                                    <Text fontSize="md">{`${chart.name}: ${chart.count}`}</Text>
                                </Box>
                            ))}
                        </Flex>
                    </Flex>
                </Box>
            </Container>
        );
    }

    const EditableControls: React.FC = () => {
        const {
            isEditing,
            getSubmitButtonProps,
            getCancelButtonProps,
            getEditButtonProps,
        } = useEditableControls();

        return isEditing ? (
            <ButtonGroup justifyContent='center' size='sm'>
                <IconButton aria-label={t('send')} icon={<CheckIcon/>} {...getSubmitButtonProps()}/>
                <IconButton aria-label={t('cancel')} icon={<CloseIcon/>} {...getCancelButtonProps()}/>
            </ButtonGroup>
        ) : (
            <Flex justifyContent='center'>
                <IconButton aria-label={t('edit')} size='sm' icon={<EditIcon/>} {...getEditButtonProps()}/>
            </Flex>
        );
    }

    const handleSelect = (selectedKeys: any[]) => {
        const anchor = document.getElementById('testSuite'+selectedKeys[0]);
        if (anchor) {
            anchor.scrollIntoView({ behavior: 'auto', block: 'center'});
        }
    };

    return (
        <ChakraProvider>
            <Box bg="white" minH="100vh">
                {/* ヘッダー */}
                <Header project_code={project_code} is_show_menu={true}/>
                <Container maxW="container.xl" py={10} pt="6rem">
                    {isLoading ? (
                        <Flex justify="center">
                            <Spinner size="xl"/>
                        </Flex>
                    ) : (
                        <>
                            <Flex justifyContent="space-around" alignItems="center" mb={4}>
                                <Chart/>
                            </Flex>
                        </>
                    )}
                    <HStack justifyContent="space-between" mb={4}>
                        <Input
                            placeholder={t('search_test_run')}
                            width="auto"
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                        {user.permissions && user.permissions.includes('edit') && (<Button colorScheme="blue" onClick={onOpen}>{t('add_new_test_run')}</Button>)}
                    </HStack>
                    <Table variant="simple">
                        <Thead>
                            <Tr>
                                <Th>{t('name')}</Th>
                                <Th>{t('test_cases_count')}</Th>
                                <Th>{t('started_at')}</Th>
                                <Th>{t('completed_at')}</Th>
                                <Th>{t('last_updated_by')}</Th>
                                {user.permissions && user.permissions.includes('edit') && (<Th>{t('actions')}</Th>)}
                            </Tr>
                        </Thead>
                        <Tbody>
                            {testRuns.filter(testRun =>
                                testRun.title.toLowerCase().includes(searchTerm.toLowerCase())
                            ).map((filteredTestRun) => (
                                <Tr key={filteredTestRun.id}>
                                    <Td><Link color="blue.500"
                                              href={`${project_code ? `/${project_code}` : ''}/runs/${filteredTestRun.id}`}>{filteredTestRun.title}</Link></Td>
                                    <Td>{filteredTestRun.count}</Td>
                                    <Td>{filteredTestRun.started_at ? filteredTestRun.started_at : "-"}</Td>
                                    <Td>{filteredTestRun.completed_at ? filteredTestRun.started_at : "-"}</Td>
                                    <Td>{filteredTestRun.updated_by.name}</Td>
                                    {user.permissions && user.permissions.includes('edit') && (
                                        <Td>
                                            <Button colorScheme="blue"
                                                    onClick={() => handleEditTestRun(filteredTestRun)}>{t('edit')}</Button>
                                        </Td>
                                    )}
                                </Tr>
                            ))}
                        </Tbody>
                    </Table>
                </Container>
            </Box>
            <MenuBar/>
            <Modal isOpen={isOpen} onClose={onClose}>
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>{t('add_new_test_run')}</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <FormControl>
                            <FormLabel>{t('title')}</FormLabel>
                            <Input value={newTestRun.title}
                                   onChange={(e) => setNewTestRun({...newTestRun, title: e.target.value})}/>
                        </FormControl>
                    </ModalBody>
                    <ModalFooter>
                        <Button colorScheme="blue" mr={3} onClick={handleAddTestRun}>{t('add')}</Button>
                        <Button variant="ghost" onClick={onClose}>{t('cancel')}</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
            <Modal isOpen={isEditModalOpen} onClose={onEditModalClose} size="full">
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>{t('edit_test_run')}</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <Editable defaultValue={selectedTestRun?.title}
                                  display="flex"
                                  submitOnBlur={false}
                                  onSubmit={handleSaveTitle}
                        >
                            <EditablePreview/>
                            <EditableInput width={300} mr={2} mb={2}
                                           onKeyDown={(e) => {
                                               if (e.key === 'Enter') {
                                                   e.preventDefault();
                                               }
                                           }}
                                           onChange={(e) => setSelectedTestRun(selectedTestRun ? {
                                               ...selectedTestRun,
                                               title: e.target.value
                                           } : null)}
                            />
                            <EditableControls/>
                        </Editable>
                        <Flex h="75vh">
                            <Box w="20%" p={5} borderRight="1px" borderColor="gray.200">
                                <Tree
                                    showLine
                                    defaultExpandAll
                                    onSelect={handleSelect}
                                    treeData={onlyTestSuites.map(suite => ({
                                        ...suite,
                                        selectable: true,
                                    }))}
                                />
                            </Box>
                            <Box w="50%" p={5} overflowY="auto" borderRight="1px" borderColor="gray.200">
                                <Box p={5}>
                                    {renderTestSuites(testSuites)}
                                </Box>
                            </Box>
                            <Box w="30%" p={5}>
                                {selectedTestCase && (
                                    <VStack align="start">
                                        <VStack align="start">
                                            <Flex justify="space-between">
                                                <Heading as="h3" size="md">{selectedTestCase.title}</Heading>
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
                                    </VStack>
                                )}
                            </Box>
                        </Flex>
                    </ModalBody>
                    <ModalFooter>
                        <Text fontSize="sm" color={isOverLimit ? "red.500" : "black"} mr={6}>
                            {`${selectedCount}/10000`}
                        </Text>
                        <Button
                            colorScheme="blue"
                            mr={3}
                            onClick={handleUpdateButtonClick}
                            isDisabled={isOverLimit}
                        >
                            {t('update')}
                        </Button>
                        <Button variant="ghost" onClick={onEditModalClose}>{t('cancel')}</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </ChakraProvider>
    );
};

export default TestPlan;
