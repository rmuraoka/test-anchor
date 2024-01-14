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
    Select,
    Table,
    Tbody,
    Td, Text,
    Th,
    Thead,
    Tr,
    useDisclosure,
    useToast
} from '@chakra-ui/react';
import {useParams} from "react-router-dom";
import Header from "../components/Header";
import {useTranslation} from "react-i18next";
import {useApiRequest} from "../components/UseApiRequest";

interface Milestone {
    id: number;
    title: string;
    description: string;
    status: string;
    due_date: string;
}

const Milestone: React.FC = () => {
    const [projectId, setProjectId] = useState<number>(0);
    const [milestones, setMilestone] = useState<Milestone[]>([]);
    const {isOpen, onOpen, onClose} = useDisclosure();
    const [newMilestone, setNewMilestone] = useState({
        project_id: projectId,
        title: '',
        status: 'Active',
        description: '',
        due_date: ''
    });
    const toast = useToast();
    const [searchTerm, setSearchTerm] = useState('');
    const {project_code} = useParams();
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [editMilestone, setEditMilestone] = useState<Milestone | null>(null);
    const openEditModal = (milestone: Milestone) => {
        setEditMilestone(milestone);
        setIsEditModalOpen(true);
    };
    const { t } = useTranslation();
    const apiRequest = useApiRequest();
    const user = JSON.parse(localStorage.getItem('user') || '{}');

    // APIからテストケースを取得
    const fetchMilestones = async () => {
        try {
            const response = await apiRequest(`/protected/${project_code}/milestones`);
            const data = await response.json();
            setProjectId(data.project_id)
            setMilestone(data.entities)
        } catch (error) {
            console.error('Error fetching Milestones:', error);
        }
    };

    useEffect(() => {
        fetchMilestones();
    }, []);

    const handleAddMilestone = async () => {
        try {
            const updatedNewMilestone = {
                ...newMilestone,
                project_id: projectId
            };
            const response = await apiRequest('/protected/milestones', {
                method: 'POST',
                body: JSON.stringify(updatedNewMilestone)
            });
            if (response.ok) {
                toast({
                    title: t('new_milestone_added'),
                    status: 'success',
                    duration: 5000,
                    isClosable: true,
                });
                setNewMilestone({
                    project_id: projectId,
                    title: '',
                    status: 'Active',
                    description: '',
                    due_date: ''
                });
                onClose();
                fetchMilestones();
            } else {
                throw new Error(t('failed_to_add_new_milestone'));
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

    const handleUpdateMilestone = async () => {
        if (editMilestone) {
            try {
                const response = await apiRequest(`/protected/milestones/${editMilestone.id}`, {
                    method: 'PUT',
                    body: JSON.stringify(editMilestone)
                });
                if (response.ok) {
                    toast({
                        title: t('milestone_updated'),
                        status: 'success',
                        duration: 5000,
                        isClosable: true,
                    });
                    setIsEditModalOpen(false);
                    fetchMilestones();
                } else {
                    throw new Error(t('milestone_update_failed'));
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

    return (
        <ChakraProvider>
            <Box bg="white" minH="100vh">
                {/* ヘッダー */}
                <Header project_code={project_code} is_show_menu={true}/>
                <Container maxW="container.xl" py={10} pt="6rem">
                    <HStack justifyContent="space-between" mb={4}>
                        <Input
                            placeholder={t('search_milestone')}
                            width="auto"
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                        {user.permissions && user.permissions.includes('edit') && (
                            <Button colorScheme="blue" onClick={onOpen}>{t('add_new_milestone')}</Button>
                        )}
                    </HStack>
                    <Table variant="simple">
                        <Thead>
                            <Tr>
                                <Th>{t('name')}</Th>
                                <Th>{t('description')}</Th>
                                <Th>{t('status')}</Th>
                                <Th>{t('due_date')}</Th>
                            </Tr>
                        </Thead>
                        <Tbody>
                            {milestones.filter(milestone =>
                                milestone.title.toLowerCase().includes(searchTerm.toLowerCase())
                            ).map((filteredMilestone) => (
                                <Tr key={filteredMilestone.id}>
                                    <Td>
                                        {user.permissions && user.permissions.includes('edit') ?
                                            <Link color="blue.500" onClick={() => openEditModal(filteredMilestone)}>
                                                {filteredMilestone.title}
                                            </Link>
                                            :
                                            <Text>{filteredMilestone.title}</Text>
                                        }
                                    </Td>
                                    <Td>{filteredMilestone.description}</Td>
                                    <Td>{filteredMilestone.status}</Td>
                                    <Td>{filteredMilestone.due_date}</Td>
                                </Tr>
                            ))}
                        </Tbody>
                    </Table>
                </Container>
            </Box>
            <Modal isOpen={isOpen} onClose={onClose}>
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>{t('add_new_milestone')}</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <FormControl>
                            <FormLabel>{t('title')}</FormLabel>
                            <Input value={newMilestone.title}
                                   onChange={(e) => setNewMilestone({...newMilestone, title: e.target.value})}/>
                            <FormLabel>{t('description')}</FormLabel>
                            <Input value={newMilestone.description}
                                   onChange={(e) => setNewMilestone({...newMilestone, description: e.target.value})}/>
                            <FormLabel>{t('due_date')}</FormLabel>
                            <Input type="date" value={newMilestone.due_date}
                                   onChange={(e) => setNewMilestone({...newMilestone, due_date: e.target.value})}/>
                        </FormControl>
                    </ModalBody>
                    <ModalFooter>
                        <Button colorScheme="blue" mr={3} onClick={handleAddMilestone}>{t('add')}</Button>
                        <Button variant="ghost" onClick={onClose}>{t('cancel')}</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
            <Modal isOpen={isEditModalOpen} onClose={() => setIsEditModalOpen(false)}>
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>{t('edit_milestone')}</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <FormControl>
                            <FormLabel>{t('title')}</FormLabel>
                            <Input
                                value={editMilestone?.title || ''}
                                onChange={(e) => setEditMilestone(editMilestone ? {
                                    ...editMilestone,
                                    title: e.target.value
                                } : null)}
                            />
                            <FormLabel>{t('description')}</FormLabel>
                            <Input
                                value={editMilestone?.description || ''}
                                onChange={(e) => setEditMilestone(editMilestone ? {
                                    ...editMilestone,
                                    description: e.target.value
                                } : null)}
                            />
                            <FormLabel>{t('due_date')}</FormLabel>
                            <Input
                                type="date"
                                value={editMilestone?.due_date || ''}
                                onChange={(e) => setEditMilestone(editMilestone ? {
                                    ...editMilestone,
                                    due_date: e.target.value
                                } : null)}
                            />
                            <FormLabel>{t('status')}</FormLabel>
                            <Select
                                value={editMilestone?.status || 'Active'}
                                onChange={(e) => setEditMilestone(editMilestone ? {
                                    ...editMilestone,
                                    status: e.target.value
                                } : null)}
                            >
                                <option value="Active">Active</option>
                                <option value="Completed">Completed</option>
                            </Select>
                        </FormControl>
                    </ModalBody>
                    <ModalFooter>
                        <Button colorScheme="blue" mr={3} onClick={handleUpdateMilestone}>{t('update')}</Button>
                        <Button variant="ghost" onClick={() => setIsEditModalOpen(false)}>{t('cancel')}</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </ChakraProvider>
    );
};

export default Milestone;
