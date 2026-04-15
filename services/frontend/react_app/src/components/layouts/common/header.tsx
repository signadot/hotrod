import { Box, Flex, Heading, HStack, Image, Text } from "@chakra-ui/react";

export const Header = () => {
    return (
        <Flex
            w='100%'
            h='80px'
            minH='80px'
            bg='gray.800'
            borderBottom='1px solid'
            borderColor='gray.700'
            px={8}
            alignItems='center'
            justifyContent='space-between'
        >
            {/* Left: HotROD logo + title */}
            <HStack spacing={4}>
                <Image src='/web_assets/hotrod_logo.png' h={16} w={16} />
                <Box>
                    <Heading size='lg' fontWeight={800} color='whiteAlpha.900' letterSpacing='-0.5px' lineHeight={1}>
                        HotROD Demo App
                    </Heading>
                    <Text fontSize='sm' color='whiteAlpha.500' mt={1}>
                        Rides On Demand · Microservices showcase
                    </Text>
                </Box>
            </HStack>

            {/* Right: Signadot branding */}
            <HStack spacing={2}>
                <Text fontSize='sm' color='whiteAlpha.500' fontWeight={500}>
                    powered by
                </Text>
                <Heading
                    size='md'
                    fontWeight={700}
                    bgGradient='linear(to-r, cyan.300, purple.300)'
                    bgClip='text'
                    letterSpacing='-0.5px'
                >
                    Signadot
                </Heading>
            </HStack>
        </Flex>
    )
}
