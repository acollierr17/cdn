import React from 'react';
import {
  Box,
  Button,
  Img,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Flex,
  Heading,
  Spacer,
  Divider,
  Center,
  Text,
} from '@chakra-ui/react';
import { Link } from 'react-router-dom';
import { BsChevronDown } from 'react-icons/bs';

interface CardProps {
  url: string;
  name: string;
}

export default function Card({ url, name }: CardProps) {
  return (
    <>
      <Box maxW="sm" borderWidth="1px" borderRadius="lg" overflow="hidden">
        <Box>
          <Center>
            <Img src={url} alt={name} />
          </Center>
        </Box>

        <Divider />

        <Flex>
          <Box p="2">
            <Heading size="md">{name}</Heading>
            <Text fontSize="xs">
              Jan 5th, 2021 &bull; 22.00Kib &bull; image/png
            </Text>
          </Box>
          <Spacer />
          <Box>
            <Menu>
              <MenuButton as={Button} rightIcon={<BsChevronDown />} m="2">
                Actions
              </MenuButton>
              <MenuList>
                <MenuItem>
                  <Link
                    to={{ pathname: url }}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    Open in New Tab
                  </Link>
                </MenuItem>
                <MenuItem>Download</MenuItem>
                <MenuItem>Delete</MenuItem>
              </MenuList>
            </Menu>
          </Box>
        </Flex>
      </Box>
    </>
  );
}
