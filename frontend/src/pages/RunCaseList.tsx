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
    Icon, IconButton,
    Input,
    Link,
    ListItem,
    Menu,
    MenuButton,
    MenuItem,
    MenuList,
    OrderedList,
    Table,
    Tbody,
    Td,
    Text,
    Textarea,
    Tr,
    UnorderedList,
    useToast,
    VStack
} from '@chakra-ui/react';
import {ChevronDownIcon, ChevronLeftIcon, ChevronRightIcon, DeleteIcon, TimeIcon} from '@chakra-ui/icons';
import {SlFolder} from "react-icons/sl";
import ReactMarkdown from 'react-markdown';
import gfm from 'remark-gfm';
import {Tree} from 'antd';
import {useNavigate, useParams} from "react-router-dom";
import Header from "../components/Header";
import {useTranslation} from "react-i18next";
import {useApiRequest} from "../components/UseApiRequest";

interface TestRunCase {
    id: number;
    testCaseId: number;
    title: string;
    content: string;
    status: Status;
    comments: Comment[];
    created_by: User;
    updated_by: User;
}

interface TestSuite {
    name: string;
    id: number;
    test_suites?: TestSuite[];
    test_cases?: TestRunCase[];
}

interface OnlyTestSuite {
    key: number;
    title: string
    children?: OnlyTestSuite[];
}

interface TestSuiteHeaderProps {
    title: string;
}

interface Comment {
    id: number;
    content: string;
    status: Status;
    created_by: User;
}

interface User {
    id: number;
    name: string;
}

interface Status {
    id: number,
    name: string,
    color: string
}

interface TestRunUpdate {
    status: string;
    started_at?: string;
    completed_at?: string | null;
    finalized_test_cases?: string;
}

const RunCaseList: React.FC = () => {
    const [testPlanId, setTestPlanId] = useState<number>(0);
    const [testRunStatus, setTestRunStatus] = useState<string>('');
    const [testRunData, setTestRunData] = useState(null);
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const [editMode, setEditMode] = useState(false); // 編集モードの状態を追加
    const [testSuites, setTestSuites] = useState<TestSuite[]>([]);
    const [onlyTestSuites, setOnlyTestSuites] = useState<OnlyTestSuite[]>([]);
    const [selectedTestCase, setSelectedTestCase] = useState<TestRunCase | null>(null);
    const [newComment, setNewComment] = useState(''); // 新しいコメントの状態
    const toast = useToast();
    const {project_code, test_run_id} = useParams();
    const {t} = useTranslation();
    const apiRequest = useApiRequest();
    const [statuses, setStatuses] = useState<Status[]>([]);

    // APIからテストケースを取得
    const fetchTestCases = async () => {
        try {
            const response = await apiRequest(`/protected/runs/${test_run_id}`);
            const data = await response.json();
            setTestRunData(data)
            setTestPlanId(data.test_plan_id);
            setTestRunStatus(data.status);
            setTestSuites(data.entities);
            setOnlyTestSuites(data.folders)
        } catch (error) {
            console.error('Error fetching TestCases:', error);
        }
    };

    const fetchStatuses = async () => {
        try {
            const response = await apiRequest(`/protected/statuses`);
            const data = await response.json();
            setStatuses(data.entities);
        } catch (error) {
            console.error('Error fetching Statuses:', error);
        }
    };

    const MenuBar = () => {
        const navigate = useNavigate();
        const handleClick = () => {
            navigate(`/${project_code}/plans/${testPlanId}`); // Replace '/new-route' with your desired path
        };
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
                <Button leftIcon={<ChevronLeftIcon/>} onClick={handleClick} mr={2}>
                    {t('back_to_test_plan')}
                </Button>
                {renderStatusButton()}
            </Flex>
        );
    };

    useEffect(() => {
        fetchTestCases();
        fetchStatuses();
    }, [project_code]);

    const handleTestCaseClick = (testCase: TestRunCase | null) => {
        setSelectedTestCase(testCase);
    };

    const handleUpdateTestCase = async () => {
        if (selectedTestCase) {
            try {
                const response = await apiRequest(`/protected/cases/${selectedTestCase.id}`, {
                    method: 'PUT',
                    body: JSON.stringify(selectedTestCase),
                });
                if (response.ok) {
                    toast({
                        title: t('test_case_updated'),
                        status: 'success',
                        duration: 5000,
                        isClosable: true,
                    });
                    setEditMode(false);
                    fetchTestCases();
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

    const handleAddComment = async () => {
        if (selectedTestCase) {
            try {
                const response = await apiRequest(`/protected/runs/cases/comments`, {
                    method: 'POST',
                    body: JSON.stringify({
                        test_run_case_id: selectedTestCase.id,
                        status_id: null,
                        content: newComment,
                        created_by_id: user.id,
                        updated_by_id: user.id,
                    }),
                });

                if (response.ok) {
                    const newCommentData = await response.json();
                    setSelectedTestCase({
                        ...selectedTestCase,
                        comments: [...selectedTestCase.comments, newCommentData]
                    });
                    fetchTestCases();

                    setNewComment(''); // コメント入力フィールドをリセット
                    toast({
                        title: t('new_comment_added'),
                        status: 'success',
                        duration: 5000,
                        isClosable: true,
                    });
                } else {
                    throw new Error(t('failed_to_add_comment'));
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

    const handleStatusChange = async (statusId: any) => {
        if (selectedTestCase) {
            try {
                const response = await apiRequest(`/protected/runs/cases/${selectedTestCase.id}`, {
                    method: 'PUT',
                    body: JSON.stringify({
                        status_id: statusId
                    })
                });

                if (!response.ok) {
                    throw new Error(t('status_update_failed'));
                }
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
        }
    };

    const handleTestRunStatusChange = async (newStatus: string, currentStatus: string) => {
        try {
            const now = new Date().toISOString();
            let updateData: TestRunUpdate = {
                status: newStatus
            };

            if (newStatus === 'InProcess' && currentStatus !== 'Completed') {
                updateData.started_at = now;
            } else if (newStatus == 'InProcess' && currentStatus == 'Completed') {
                updateData.completed_at = null;
                updateData.finalized_test_cases = '';
            } else if (newStatus === 'Completed') {
                updateData.completed_at = now;
                updateData.finalized_test_cases = JSON.stringify(testRunData);
                setEditMode(false)
            }

            const response = await apiRequest(`/protected/runs/${test_run_id}`, {
                method: 'PUT',
                body: JSON.stringify(updateData)
            });

            if (!response.ok) {
                throw new Error(t('status_update_failed'));
            }
            setTestRunStatus(newStatus);
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


    const renderStatusButton = () => {
        switch (testRunStatus) {
            case 'NotExecuted':
                return <Button leftIcon={<TimeIcon/>}
                               onClick={() => handleTestRunStatusChange('InProcess', testRunStatus)}>{t('start_test')}</Button>;
            case 'InProcess':
                return <Button leftIcon={<TimeIcon/>}
                               onClick={() => handleTestRunStatusChange('Completed', testRunStatus)}>{t('complete')}</Button>;
            case 'Completed':
                return <Button leftIcon={<TimeIcon/>}
                               onClick={() => handleTestRunStatusChange('InProcess', testRunStatus)}>{t('return_to_in_process')}</Button>;
            default:
                return null;
        }
    };

    const renderTestCases = (testRunCases: TestRunCase[]) => (
        <Table variant="simple">
            <Tbody>
                {testRunCases.map(testRunCase => (
                    <Tr cursor="pointer" _hover={{bg: "gray.100"}}>
                        <Td borderBottom="1px" borderColor="gray.200">{testRunCase.title}</Td>
                        <Td width="120px">
                            <Flex justifyContent="flex-end" alignItems="center">
                                <Menu>
                                    <MenuButton as={Button} rightIcon={<ChevronDownIcon/>}
                                                bg={`${testRunCase.status.color}.200`} width="100%">
                                        {testRunCase.status.name}
                                    </MenuButton>
                                    <MenuList>
                                        {renderStatusOptions()}
                                    </MenuList>
                                </Menu>
                                <Box mx={2}>
                                    {selectedTestCase && selectedTestCase.id === testRunCase.id ?
                                        <IconButton
                                            aria-label={t('open_test_case')}
                                            icon={<ChevronLeftIcon/>}
                                            size="sm"
                                            onClick={() => handleTestCaseClick(null)}
                                        /> :
                                        <IconButton
                                            aria-label={t('close_test_case')}
                                            icon={<ChevronRightIcon/>}
                                            size="sm"
                                            onClick={() => handleTestCaseClick(testRunCase)}
                                        />}
                                </Box>
                            </Flex>
                        </Td>
                    </Tr>
                ))}
            </Tbody>
        </Table>
    );

    const TestSuiteHeader: React.FC<TestSuiteHeaderProps> = ({title}) => {
        return (
            <Flex justifyContent="space-between" alignItems="center" mb={4}>
                <Flex alignItems="center">
                    <Icon as={SlFolder} mr={2}/>
                    <Text fontSize="lg" fontWeight="bold">{title}</Text>
                </Flex>
            </Flex>
        );
    };

    const renderTestSuites = (suites: TestSuite[], level = 0) => (
        suites.map(suite => (
            <Box key={suite.name} mb={4} pl={`${level * 2}em`}>
                <TestSuiteHeader title={suite.name}/>
                <Box mb={4}>
                    {suite.test_cases && suite.test_cases.length > 0 && renderTestCases(suite.test_cases)}
                </Box>
                {suite.test_suites && suite.test_suites.length > 0 && renderTestSuites(suite.test_suites, level + 1)}
            </Box>
        ))
    );

    const renderStatusOptions = () => {
        return statuses.map(status => (
            <MenuItem
                key={status.id}
                value={status.id}
                bg={`${status.color}.200`}
                onClick={() => handleStatusChange(status.id)}
            >
                {status.name}
            </MenuItem>
        ));
    };

    const renderComments = (comments: any[]) => {
        return comments.map((comment, index) => (
            <Box key={index} mb={2} p={4} border="1px" borderColor="gray.200" borderRadius="md" w="90%">
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
                    {comment.content}
                </ReactMarkdown>
                <Flex justify="space-between" mt={2}>
                    <Text fontSize="sm">by {comment.created_by.name}</Text>
                </Flex>
            </Box>
        ));
    };

    return (
        <ChakraProvider>
            <Header project_code={project_code} is_show_menu={true}/>
            <Flex h="100vh">
                <Box w="20%" p={5} borderRight="1px" borderColor="gray.200" pt="6rem" overflowY="auto">
                    <Tree
                        showLine
                        defaultExpandAll
                        treeData={onlyTestSuites}
                    />
                </Box>
                <Box w="50%" p={5} overflowY="auto" borderRight="1px" borderColor="gray.200" pt="6rem">
                    <Box p={5}>
                        {renderTestSuites(testSuites)}
                    </Box>
                </Box>
                <Box w="30%" p={5} pt="6rem" overflowY="auto">
                    {selectedTestCase && (
                        <VStack align="start">
                            {editMode ? (
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
                                        <FormLabel> {t('content')}</FormLabel>
                                        <Textarea
                                            height={200}
                                            value={selectedTestCase?.content}
                                            onChange={(e) => setSelectedTestCase({
                                                ...selectedTestCase,
                                                content: e.target.value
                                            })}
                                        />
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
                                            {testRunStatus !== 'Completed' && (
                                                <Button
                                                    size="sm"
                                                    onClick={() => setEditMode(true)}
                                                    ml={2}
                                                >
                                                    {t('edit')}
                                                </Button>
                                            )}
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
                                    <Heading as="h4" size="md" mt={6} mb={2}>{t('comments')}</Heading>
                                    {renderComments(selectedTestCase.comments)}
                                    <Flex p={5} borderTop="1px" borderColor="gray.200">
                                        <Textarea
                                            placeholder={t('add_comment')}
                                            value={newComment}
                                            onChange={(e) => setNewComment(e.target.value)}
                                            mr={2}
                                        />
                                        <Button colorScheme="blue" onClick={handleAddComment}>{t('send')}</Button>
                                    </Flex>
                                </>
                            )}
                        </VStack>
                    )}
                </Box>
            </Flex>
            <MenuBar/>
        </ChakraProvider>
    );
};

export default RunCaseList;
