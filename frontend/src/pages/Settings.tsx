import React from 'react';
import {Box, ChakraProvider, Flex} from '@chakra-ui/react';
import Header from "../components/Header";
import SettingMenu from "../components/SettingMenu";
import {useTranslation} from "react-i18next";

const Settings: React.FC = () => {
    const { t } = useTranslation();

    return (
        <ChakraProvider>
            <Flex direction="column">
                <Header project_code={undefined} is_show_menu={false}/>
                <Flex>
                    <SettingMenu/>
                    <Box p={8} pt="6rem" w="80%">
                        {t('select_menu')}
                    </Box>
                </Flex>
            </Flex>
        </ChakraProvider>
    );
};

export default Settings;
