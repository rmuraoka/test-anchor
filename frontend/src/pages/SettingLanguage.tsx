import React, {useState} from 'react';
import {Box, Button, ChakraProvider, Flex, FormControl, FormLabel, Input, Select, useToast} from '@chakra-ui/react';
import Header from "../components/Header";
import SettingMenu from "../components/SettingMenu";
import {useTranslation} from "react-i18next";
import {useApiRequest} from "../components/UseApiRequest";
import i18n from "i18next";

const SettingLanguage: React.FC = () => {
    const {t} = useTranslation();
    const [language, setLanguage] = useState('');
    const toast = useToast();
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const apiRequest = useApiRequest();

    const handleSubmit = async () => {
        try {
            const response = await apiRequest(`/protected/members/${user.id}`, {
                method: 'PUT',
                body: JSON.stringify({language: language})
            });

            const storedUser = localStorage.getItem('user');
            if (storedUser) {
                const user = JSON.parse(storedUser);
                user.language = language;
                localStorage.setItem('user', JSON.stringify(user));
                i18n.changeLanguage(language);
            }

            toast({title: t('language_changed_successfully'), status: 'success'});
        } catch (error) {
            toast({title: t('failed_to_change_language'), status: 'error'});
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
                            <FormLabel>{t('select_language')}</FormLabel>
                            <Select
                                width={300}
                                placeholder={t('select_language')}
                                onChange={(e) => setLanguage(e.target.value)}
                                value={language}
                            >
                                <option value="en">English</option>
                                <option value="ja">日本語</option>
                            </Select>
                            <Button mt={4} colorScheme="blue" onClick={handleSubmit}>{t('change')}</Button>
                        </FormControl>
                    </Box>
                </Flex>
            </Flex>
        </ChakraProvider>
    );
};

export default SettingLanguage;
