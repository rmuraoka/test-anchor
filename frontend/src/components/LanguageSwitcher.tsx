import React, { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { IconButton, Menu, MenuButton, MenuItem, MenuList } from '@chakra-ui/react';
import { AiOutlineGlobal } from "react-icons/ai";

function LanguageSwitcher() {
    const { i18n } = useTranslation();

    useEffect(() => {
        const storedUser = localStorage.getItem('user');
        if (storedUser) {
            const user = JSON.parse(storedUser);
            if (user.language) {
                i18n.changeLanguage(user.language);
            }
        }
    }, [i18n]);

    const changeLanguage = (languageCode: string | undefined) => {
        const storedUser = localStorage.getItem('user');
        if (storedUser) {
            const user = JSON.parse(storedUser);
            user.language = languageCode;
            localStorage.setItem('user', JSON.stringify(user));
            i18n.changeLanguage(languageCode);
        }
    };

    return (
        <Menu>
            <MenuButton as={IconButton} icon={<AiOutlineGlobal />} variant="outline" />
            <MenuList>
                <MenuItem onClick={() => changeLanguage('en')}>English</MenuItem>
                <MenuItem onClick={() => changeLanguage('ja')}>日本語</MenuItem>
            </MenuList>
        </Menu>
    );
}

export default LanguageSwitcher;
