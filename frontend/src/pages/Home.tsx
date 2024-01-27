import React, {useEffect, useState} from 'react';
import {
    Box,
    Button,
    ChakraProvider,
    Container,
    FormControl, FormHelperText,
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
    useToast
} from '@chakra-ui/react';
import Header from "../components/Header";
import {useTranslation} from 'react-i18next';
import {useApiRequest} from "../components/UseApiRequest";

interface Project {
    id: number;
    title: string;
    code: string;
    description: string;
}

const Home: React.FC = () => {
    const [projects, setProjects] = useState<Project[]>([]);
    const {isOpen, onOpen, onClose} = useDisclosure();
    const [newProject, setNewProject] = useState({title: '', code: '', description: ''});
    const [searchTerm, setSearchTerm] = useState('');
    const [isSubmitDisabled, setIsSubmitDisabled] = useState(true);
    const toast = useToast();
    const apiRequest = useApiRequest();
    const user = JSON.parse(localStorage.getItem('user') || '{}');

    const fetchProjects = async () => {
        try {
            const response = await apiRequest(`/protected/projects`);
            const data = await response.json();
            setProjects(data.entities);
        } catch (error) {
            console.error('Error fetching projects:', error);
        }
    };
    const {t} = useTranslation();

    useEffect(() => {
        fetchProjects();
    }, []);

    const handleAddProject = async () => {
        try {
            const response = await apiRequest('/protected/projects', {
                method: 'POST',
                body: JSON.stringify(newProject)
            });

            if (response.ok) {
                toast({
                    title: t('new_project_added'),
                    status: 'success',
                    duration: 5000,
                    isClosable: true,
                });
                setNewProject({
                    title: '',
                    code: '',
                    description: '',
                });
                onClose();
                setIsSubmitDisabled(true);
                fetchProjects();
            } else {
                throw new Error(t('failure_add_project'));
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

    const validateForm = () => {
        const isFormValid = newProject.title && newProject.code && newProject.description;
        setIsSubmitDisabled(!isFormValid);
    };

    return (
        <ChakraProvider>
            <Header is_show_menu={false} project_code={undefined}/>
            <Box bg="white" minH="100vh" pt="4rem">
                <Container maxW="container.xl" py={10}>
                    <HStack justifyContent="space-between" mb={4}>
                        <Input
                            placeholder={t('search_project')}
                            width="auto"
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                        {user.permissions && user.permissions.includes('edit') && (
                            <Button colorScheme="blue" onClick={onOpen}>{t('add_project')}</Button>
                        )}
                    </HStack>
                    <Table variant="simple">
                        <Thead>
                            <Tr>
                                <Th>{t('project_name')}</Th>
                                <Th>{t('code')}</Th>
                                <Th>{t('description')}</Th>
                            </Tr>
                        </Thead>
                        <Tbody>
                            {projects.filter(projects =>
                                projects.title.toLowerCase().includes(searchTerm.toLowerCase())
                            ).map((filteredProject) => (
                                <Tr key={filteredProject.id}>
                                    <Td><Link color="blue.500"
                                              href={`${filteredProject.code}`}>{filteredProject.title}</Link></Td>
                                    <Td>{filteredProject.code}</Td>
                                    <Td>{filteredProject.description}</Td>
                                </Tr>
                            ))}
                        </Tbody>
                    </Table>
                </Container>
            </Box>
            <Modal isOpen={isOpen} onClose={onClose}>
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>{t('add_project')}</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <FormControl isRequired>
                            <FormLabel>{t('project_name')}</FormLabel>
                            <Input data-test="input title"
                                   value={newProject.title}
                                   onChange={(e) => {
                                       setNewProject({...newProject, title: e.target.value});
                                       validateForm();
                                   }}
                                   onBlur={() => validateForm()}
                            />
                            <FormLabel mt={2}>{t('code')}</FormLabel>
                            <FormHelperText color="gray.500">{t('code_help_text')}</FormHelperText>
                            <Input data-test="input code"
                                   value={newProject.code}
                                   onChange={(e) => {
                                       const inputValue = e.target.value;
                                       const isValidInput = /^[A-Z0-9]*$/.test(inputValue);

                                       if (isValidInput) {
                                           setNewProject({...newProject, code: inputValue.toUpperCase()});
                                           validateForm();
                                       }
                                   }}
                                   onKeyDown={(e: React.KeyboardEvent<HTMLInputElement>) => {
                                       const isAlphabetKey = /^[a-z]$/.test(e.key);
                                       if (isAlphabetKey) {
                                           const uppercaseInputValue = e.key.toUpperCase();
                                           const newValue = (e.target as HTMLInputElement).value + uppercaseInputValue;
                                           setNewProject({...newProject, code: newValue});
                                           validateForm();
                                           return;
                                       }
                                   }}
                                   onBlur={() => validateForm()}
                            />
                            <FormLabel mt={2}>{t('description')}</FormLabel>
                            <Input data-test="input description"
                                   value={newProject.description}
                                   onChange={(e) => {
                                       setNewProject({...newProject, description: e.target.value});
                                       validateForm();
                                   }}
                                   onBlur={() => validateForm()}
                            />
                        </FormControl>
                    </ModalBody>
                    <ModalFooter>
                        <Button colorScheme="blue"
                                data-test="button project add"
                                mr={3}
                                onClick={handleAddProject}
                                isDisabled={isSubmitDisabled}>{t('add')}</Button>
                        <Button variant="ghost" onClick={onClose}>{t('cancel')}</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </ChakraProvider>
    );
};

export default Home;
