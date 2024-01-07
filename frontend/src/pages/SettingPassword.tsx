import React, {useState} from 'react';
import {Box, Button, ChakraProvider, Flex, FormControl, FormLabel, Input, useToast} from '@chakra-ui/react';
import Header from "../components/Header";
import SettingMenu from "../components/SettingMenu";
import {useTranslation} from "react-i18next";
import {useApiRequest} from "../components/UseApiRequest";

const SettingPassword: React.FC = () => {
    const {t} = useTranslation();
    const [currentPassword, setCurrentPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const toast = useToast();
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const apiRequest = useApiRequest();

    const handleSubmit = async () => {
        // バリデーションチェック
        if (newPassword !== confirmPassword) {
            toast({title: t('passwords_do_not_match'), status: 'error'});
            return;
        }

        // APIリクエスト
        try {
            const response = await apiRequest(`/protected/members/${user.id}`, {
                method: 'PUT',
                body: JSON.stringify({current_password: currentPassword, password: newPassword})
            });

            toast({title: t('password_changed_successfully'), status: 'success'});
        } catch (error) {
            toast({title: t('error_updating_password'), status: 'error'});
        }
    };

    return (
        <ChakraProvider>
            <Flex direction="column">
                <Header project_code={undefined} is_show_menu={false}/>
                <Flex>
                    <SettingMenu/>
                    <Box p={8} pt="6rem" w="80%">
                        <FormControl>
                            <Box>
                                <FormLabel>{t('current_password')}</FormLabel>
                                <Input w={400} mb={2} type="password" value={currentPassword}
                                       onChange={e => setCurrentPassword(e.target.value)}/>
                                <FormLabel>{t('new_password')}</FormLabel>
                                <Input w={400} mb={2} type="password" value={newPassword}
                                       onChange={e => setNewPassword(e.target.value)}/>
                                <FormLabel>{t('confirm_new_password')}</FormLabel>
                                <Input w={400} mb={2} type="password" value={confirmPassword}
                                       onChange={e => setConfirmPassword(e.target.value)}/>
                            </Box>
                            <Button mt={4} colorScheme="blue" onClick={handleSubmit}>{t('change_password')}</Button>
                        </FormControl>
                    </Box>
                </Flex>
            </Flex>
        </ChakraProvider>
    );
};

export default SettingPassword;
