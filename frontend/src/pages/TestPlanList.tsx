import React, {useEffect, useState} from 'react';
import {
    Box,
    Button,
    ChakraProvider,
    Container,
    FormControl,
    FormLabel,
    HStack,
    Input,
    Link,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    Table,
    Tbody,
    Td,
    Th,
    Thead,
    Tr,
    useDisclosure,
    useToast,
} from '@chakra-ui/react';
import {useParams} from "react-router-dom";
import Header from "../components/Header";
import {useTranslation} from "react-i18next";
import {useApiRequest} from "../components/UseApiRequest";

interface TestPlan {
    id: number;
    title: string;
    status: string;
    started_at: string | null;
    completed_at: string | null;
    created_by: User;
    updated_by: User;
}

interface User {
    id: number;
    name: string;
}

const TestPlanList: React.FC = () => {
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const [projectId, setProjectId] = useState<number>(0);
    const [testPlans, setTestPlan] = useState<TestPlan[]>([]);
    const {isOpen, onOpen, onClose} = useDisclosure();
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [newTestPlan, setNewTestPlan] = useState({
        project_id: projectId,
        title: '',
        status: 'NotExecuted',
        created_by_id: user.id,
        updated_by_id: user.id
    });
    const [editTestPlan, setEditTestPlan] = useState<TestPlan | null>(null);
    const openEditModal = (testPlan: TestPlan) => {
        setEditTestPlan(testPlan);
        setIsEditModalOpen(true);
    };
    const [searchTerm, setSearchTerm] = useState('');
    const toast = useToast();
    const {project_code} = useParams();
    const {t} = useTranslation();
    const apiRequest = useApiRequest();

    const fetchTestPlans = async () => {
        try {
            const response = await apiRequest(`/protected/${project_code}/plans`);
            const data = await response.json();
            setProjectId(data.project_id)
            setTestPlan(data.entities)
        } catch (error) {
            console.error('Error fetching TestPlans:', error);
        }
    };

    useEffect(() => {
        fetchTestPlans();
    }, []);

    const handleAddTestPlan = async () => {
        try {
            const updatedNewTestPlan = {
                ...newTestPlan,
                project_id: projectId
            };
            const response = await apiRequest(`/protected/plans`, {
                method: 'POST',
                body: JSON.stringify(updatedNewTestPlan)
            });
            if (response.ok) {
                toast({
                    title: t('new_test_plan_added'),
                    status: 'success',
                    duration: 5000,
                    isClosable: true,
                });
                setNewTestPlan({
                    project_id: projectId,
                    title: '',
                    status: 'NotExecuted',
                    created_by_id: user.id,
                    updated_by_id: user.id
                });
                onClose();
                fetchTestPlans();
            } else {
                throw new Error(t('failed_to_add_new_test_plan'));
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

    const handleEditTestPlan = async () => {
        try {
            if (!editTestPlan) {
                return;
            }
            const response = await apiRequest(`/protected/plans/${editTestPlan.id}`, {
                method: 'PUT',
                body: JSON.stringify({title: editTestPlan.title, updated_by_id: user.id})
            });
            if (response.ok) {
                toast({
                    title: t('test_plan_updated'),
                    status: 'success',
                    duration: 5000,
                    isClosable: true,
                });
                setNewTestPlan({
                    project_id: projectId,
                    title: '',
                    status: 'NotExecuted',
                    created_by_id: user.id,
                    updated_by_id: user.id
                });
                setIsEditModalOpen(false);
                fetchTestPlans();
            } else {
                throw new Error(t('failed_to_update_test_plan'));
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

    return (
        <ChakraProvider>
            <Box bg="white" minH="100vh">
                {/* ヘッダー */}
                <Header project_code={project_code} is_show_menu={true}/>
                <Container maxW="container.xl" py={10} pt="6rem">
                    <HStack justifyContent="space-between" mb={4}>
                        <Input
                            placeholder={t('search_test_plan')}
                            width="auto"
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                        <Button colorScheme="blue" onClick={onOpen}>{t('add_new_test_plan')}</Button>
                    </HStack>
                    <Table variant="simple">
                        <Thead>
                            <Tr>
                                <Th>{t('name')}</Th>
                                <Th>{t('status')}</Th>
                                <Th>{t('started_at')}</Th>
                                <Th>{t('completed_at')}</Th>
                                <Th>{t('last_updated_by')}</Th>
                                <Th>{t('action')}</Th>
                            </Tr>
                        </Thead>
                        <Tbody>
                            {testPlans.filter(testPlan =>
                                testPlan.title.toLowerCase().includes(searchTerm.toLowerCase())
                            ).map((filteredTestPlan) => (
                                <Tr key={filteredTestPlan.id}>
                                    <Td>
                                        <Link color="blue.500"
                                              href={`${project_code ? `/${project_code}` : ''}/plans/${filteredTestPlan.id}`}>
                                            {filteredTestPlan.title}
                                        </Link>
                                    </Td>
                                    <Td>{filteredTestPlan.status}</Td>
                                    <Td>{filteredTestPlan.started_at ? filteredTestPlan.started_at : "-"}</Td>
                                    <Td>{filteredTestPlan.completed_at ? filteredTestPlan.started_at : "-"}</Td>
                                    <Td>{filteredTestPlan.updated_by.name}</Td>
                                    <Td><Button colorScheme="blue"
                                                onClick={() => openEditModal(filteredTestPlan)}>{t('edit')}</Button></Td>
                                </Tr>
                            ))}
                        </Tbody>
                    </Table>
                </Container>
            </Box>
            <Modal isOpen={isOpen} onClose={onClose}>
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>{t('add_new_test_plan')}</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <FormControl>
                            <FormLabel>{t('title')}</FormLabel>
                            <Input value={newTestPlan.title}
                                   onChange={(e) => setNewTestPlan({...newTestPlan, title: e.target.value})}/>
                        </FormControl>
                    </ModalBody>
                    <ModalFooter>
                        <Button colorScheme="blue" mr={3} onClick={handleAddTestPlan}>{t('add')}</Button>
                        <Button variant="ghost" onClick={onClose}>{t('cancel')}</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
            <Modal isOpen={isEditModalOpen} onClose={() => setIsEditModalOpen(false)}>
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>{t('edit_test_plan')}</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <FormControl>
                            <FormLabel>{t('title')}</FormLabel>
                            <Input value={editTestPlan?.title || ''}
                                   onChange={(e) => setEditTestPlan(editTestPlan ? {
                                       ...editTestPlan,
                                       title: e.target.value
                                   } : null)}
                            />
                        </FormControl>
                    </ModalBody>
                    <ModalFooter>
                        <Button colorScheme="blue" mr={3} onClick={handleEditTestPlan}>{t('update')}</Button>
                        <Button variant="ghost" onClick={() => setIsEditModalOpen(false)}>{t('cancel')}</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </ChakraProvider>
    );
};

export default TestPlanList;
