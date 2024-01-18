import React from "react";
import {Box, Divider, Flex, Icon, Link, Text, VStack} from "@chakra-ui/react";
import {FaUser} from "react-icons/fa";
import {LockIcon} from "@chakra-ui/icons";
import {useTranslation} from "react-i18next";
import {AiOutlineGlobal} from "react-icons/ai";

const Header: React.FC = () => {
    const {t} = useTranslation();
    const user = JSON.parse(localStorage.getItem('user') || '{}');

    return (
        <Box w="20%" bg="gray.100" p={4} h="100vh" pt="6rem">
            <VStack align="stretch" spacing={4}>
                {user.permissions && user.permissions.includes('admin') && (
                    <>
                        <Link href="/settings/members">
                            <Flex align="center">
                                <Icon as={FaUser} mr={2}/>
                                <Text fontSize="lg" fontWeight="bold">{t('member_management')}</Text>
                            </Flex>
                        </Link>
                        <Divider borderColor='gray.400'/>
                    </>)
                }
                <Link href="/settings/password">
                    <Flex align="center">
                        <LockIcon mr={2}/>
                        <Text fontSize="lg" fontWeight="bold">{t('change_password')}</Text>
                    </Flex>
                </Link>
                <Link href="/settings/language">
                    <Flex align="center">
                        <AiOutlineGlobal/>
                        <Text ml={2} fontSize="lg" fontWeight="bold">{t('language')}</Text>
                    </Flex>
                </Link>
            </VStack>
        </Box>
    );
};

export default Header;
