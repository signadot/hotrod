import { Divider, Flex, Heading, HStack, Image, Text } from "@chakra-ui/react";

export const Header = () => {
    return (
        <Flex
            w='100%'
            h='56px'
            minH='56px'
            bg='gray.800'
            borderBottom='1px solid'
            borderColor='gray.700'
            px={6}
            alignItems='center'
        >
            <HStack spacing={3}>
                <Image src='/web_assets/hotrod_logo.png' h={9} w={9} />
                <Heading size='md' fontWeight={700} color='whiteAlpha.900'>
                    HotROD
                </Heading>
                <Divider orientation='vertical' h='24px' borderColor='gray.600' />
                <Text fontSize='sm' color='whiteAlpha.500'>
                    by Signadot
                </Text>
            </HStack>
        </Flex>
    )
}