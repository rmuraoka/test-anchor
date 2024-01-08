import React, {useEffect, useState} from 'react';
import {
    Box,
    Button,
    ChakraProvider,
    Flex,
    FormControl,
    FormLabel,
    HStack,
    IconButton,
    Input,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    Table,
    Tbody,
    Td, Text,
    Th,
    Thead,
    Tr,
    useToast
} from '@chakra-ui/react';
import Header from "../components/Header";
import SettingMenu from "../components/SettingMenu";
import {EditIcon, EmailIcon, LockIcon, NotAllowedIcon, UnlockIcon} from "@chakra-ui/icons";
import {useTranslation} from "react-i18next";
import {useApiRequest} from "../components/UseApiRequest";

interface Member {
    id: number;
    email: string;
    name: string;
    status: string;
}

const SettingMembers: React.FC = () => {
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const [searchTerm, setSearchTerm] = useState('');
    const [members, setMembers] = useState<Member[]>([]);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [isConfirmResendModalOpen, setIsConfirmResendModalOpen] = useState(false);
    const [isConfirmLockModalOpen, setIsConfirmLockModalOpen] = useState(false);
    const [isConfirmUnLockModalOpen, setIsConfirmUnLockModalOpen] = useState(false);
    const [targetMemberId, setTargetMemberId] = useState<number|null>(null);
    const [editableMember, setEditableMember] = useState({id: 0, name: '', email: ''});
    const [isInviteModalOpen, setIsInviteModalOpen] = useState(false);
    const [newMember, setNewMember] = useState({name: '', email: ''});
    const toast = useToast();
    const { t } = useTranslation();
    const apiRequest = useApiRequest();

    const fetchMembers = async () => {
        try {
            const response = await apiRequest(`/protected/members`);
            const data = await response.json();
            setMembers(data.entities)
        } catch (error) {
            console.error('Error fetching TestCases:', error);
        }
    };

    const handleInvite = async () => {
        try {
            const response = await apiRequest(`/protected/members`, {
                method: 'POST',
                body: JSON.stringify(newMember)
            });

            if (!response.ok) {
                throw new Error('Failed to invite member');
            }

            setIsInviteModalOpen(false);
            setNewMember({name: '', email: ''});
            toast({
                title: t('member_invited_successfully'),
                status: 'success',
                duration: 5000,
                isClosable: true
            });
            fetchMembers(); // Reload the member list
        } catch (error) {
            console.error('Error inviting member:', error);
        }
    };

    const handleUpdateMember = async () => {
        try {
            const response = await apiRequest(`/protected/members/${editableMember.id}`, {
                method: 'PUT',
                body: JSON.stringify(editableMember)
            });
            setIsEditModalOpen(false);
            fetchMembers();
        } catch (error) {
            console.error('Error updating member:', error);
        }
    };

    const handleResendPassword = async () => {
        try {
            const response = await apiRequest(`/protected/members/${targetMemberId}?resend_password=true`, {
                method: 'PUT',
                body: JSON.stringify({})
            });
            setIsConfirmResendModalOpen(false);
            fetchMembers();
        } catch (error) {
            console.error('Error updating member:', error);
        }
    };

    const handleUpdateDeactivate = async () => {
        try {
            const response = await apiRequest(`/protected/members/${targetMemberId}`, {
                method: 'PUT',
                body: JSON.stringify({status: "Deactivate"})
            });
            setIsConfirmLockModalOpen(false);
            fetchMembers();
        } catch (error) {
            console.error('Error updating member:', error);
        }
    };

    const handleUpdateActivate = async () => {
        try {
            const response = await apiRequest(`/protected/members/${targetMemberId}`, {
                method: 'PUT',
                body: JSON.stringify({status: "Active"})
            });
            setIsConfirmUnLockModalOpen(false);
            fetchMembers();
        } catch (error) {
            console.error('Error updating member:', error);
        }
    };

    const handleResendEmail = (memberId: number) => {
        setTargetMemberId(memberId);
        setIsConfirmResendModalOpen(true);
    };

    useEffect(() => {
        fetchMembers();
    }, []);

    const handleEditMember = (member: Member) => {
        setEditableMember({id: member.id, name: member.name, email: member.email});
        setIsEditModalOpen(true);
    };

    const handleDeactivateMember = (memberId: number) => {
        setTargetMemberId(memberId);
        setIsConfirmLockModalOpen(true);
    };

    const handleActivateMember = (memberId: number) => {
        setTargetMemberId(memberId);
        setIsConfirmUnLockModalOpen(true);
    };

    return (
        <ChakraProvider>
            <Flex direction="column">
                <Header project_code={undefined} is_show_menu={false}/>
                <Flex>
                    <SettingMenu/>
                    <Box p={8} pt="6rem" w="80%">
                        <HStack justifyContent="space-between" mb={4}>
                            <Input
                                placeholder={t('search_member')}
                                width="auto"
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                            />
                            <Button colorScheme="blue" onClick={() => setIsInviteModalOpen(true)}>{t('invite')}</Button>
                        </HStack>
                        <Table variant="simple">
                            <Thead>
                                <Tr>
                                    <Th>{t('name')}</Th>
                                    <Th>{t('email')}</Th>
                                    <Th>{t('actions')}</Th>
                                </Tr>
                            </Thead>
                            <Tbody>
                                {members.filter(members =>
                                    members.name.toLowerCase().includes(searchTerm.toLowerCase())
                                ).map((filteredMember) => (
                                    <Tr key={filteredMember.id}>
                                        <Td>{filteredMember.name}</Td>
                                        <Td>{filteredMember.email}</Td>
                                        <Td>
                                            <IconButton
                                                aria-label={t('edit_member')}
                                                icon={<EditIcon/>}
                                                mr={2}
                                                onClick={() => handleEditMember(filteredMember)}
                                            />
                                            {user.id !== filteredMember.id && (
                                                <IconButton
                                                    aria-label={t('resend_email')}
                                                    icon={<EmailIcon/>}
                                                    mr={2}
                                                    colorScheme="blue"
                                                    onClick={() => handleResendEmail(filteredMember.id)}
                                                />
                                            )}
                                            {user.id !== filteredMember.id && filteredMember.status === "Active" && (
                                                <IconButton
                                                    aria-label={t('deactivate_member')}
                                                    icon={<NotAllowedIcon/>}
                                                    colorScheme="red"
                                                    onClick={() => handleDeactivateMember(filteredMember.id)}
                                                />
                                            )}
                                            {user.id !== filteredMember.id && filteredMember.status === "Deactivate" && (
                                                <IconButton
                                                    aria-label={t('deactivate_member')}
                                                    icon={<UnlockIcon/>}
                                                    colorScheme="blue"
                                                    onClick={() => handleActivateMember(filteredMember.id)}
                                                />
                                            )}
                                        </Td>
                                    </Tr>
                                ))}
                            </Tbody>
                        </Table>
                    </Box>
                </Flex>
            </Flex>
            <Modal isOpen={isEditModalOpen} onClose={() => setIsEditModalOpen(false)}>
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>{t('edit_member')}</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <FormControl>
                            <FormLabel>{t('name')}</FormLabel>
                            <Input
                                value={editableMember.name}
                                onChange={(e) => setEditableMember({...editableMember, name: e.target.value})}
                            />
                            <FormLabel>{t('email')}</FormLabel>
                            <Input
                                type="email"
                                value={editableMember.email}
                                onChange={(e) => setEditableMember({...editableMember, email: e.target.value})}
                            />
                        </FormControl>
                    </ModalBody>
                    <ModalFooter>
                        <Button colorScheme="blue" mr={3} onClick={handleUpdateMember}>{t('dupdate')}</Button>
                        <Button variant="ghost" onClick={() => setIsEditModalOpen(false)}>{t('cancel')}</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
            <Modal isOpen={isInviteModalOpen} onClose={() => setIsInviteModalOpen(false)}>
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>{t('invite_member')}</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <FormControl>
                            <FormLabel>{t('name')}</FormLabel>
                            <Input value={newMember.name}
                                   onChange={(e) => setNewMember({...newMember, name: e.target.value})}/>
                            <FormLabel>{t('email')}</FormLabel>
                            <Input type="email" value={newMember.email}
                                   onChange={(e) => setNewMember({...newMember, email: e.target.value})}/>
                        </FormControl>
                    </ModalBody>
                    <ModalFooter>
                        <Button colorScheme="blue" mr={3} onClick={handleInvite}>{t('invite')}</Button>
                        <Button variant="ghost" onClick={() => setIsInviteModalOpen(false)}>{t('cancel')}</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
            <Modal isOpen={isConfirmLockModalOpen} onClose={() => setIsConfirmLockModalOpen(false)}>
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>{t('deactivate_member')}</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <Text>{t('confirm_deactivate')}</Text>
                    </ModalBody>
                    <ModalFooter>
                        <Button colorScheme="red" onClick={() => {handleUpdateDeactivate()}}>
                            {t('deactivate')}
                        </Button>
                        <Button ml={2} onClick={() => setIsConfirmLockModalOpen(false)}>{t('cancel')}</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
            <Modal isOpen={isConfirmUnLockModalOpen} onClose={() => setIsConfirmUnLockModalOpen(false)}>
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>{t('activate_member')}</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <Text>{t('confirm_activate')}</Text>
                    </ModalBody>
                    <ModalFooter>
                        <Button colorScheme="blue" onClick={() => {handleUpdateActivate()}}>
                            {t('activate')}
                        </Button>
                        <Button ml={2} onClick={() => setIsConfirmUnLockModalOpen(false)}>{t('cancel')}</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
            <Modal isOpen={isConfirmResendModalOpen} onClose={() => setIsConfirmResendModalOpen(false)}>
                <ModalOverlay/>
                <ModalContent>
                    <ModalHeader>{t('resend_email')}</ModalHeader>
                    <ModalCloseButton/>
                    <ModalBody>
                        <Text>{t('confirm_resend_email')}</Text>
                    </ModalBody>
                    <ModalFooter>
                        <Button colorScheme="blue" onClick={() => {handleResendPassword()}}>
                            {t('send')}
                        </Button>
                        <Button ml={2} onClick={() => setIsConfirmResendModalOpen(false)}>{t('cancel')}</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </ChakraProvider>
    );
};

export default SettingMembers;
