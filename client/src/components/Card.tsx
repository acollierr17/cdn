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
  Skeleton,
  SkeletonText,
} from '@chakra-ui/react';
import { Link as RouterLink } from 'react-router-dom';
import { BsChevronDown } from 'react-icons/bs';
import type { ImageResult } from '../api';
import { formatDate, formatFileSize } from '../functions';

interface CardProps extends ImageResult {
  loading: boolean;
  token: string;
}

export default function Card(props: CardProps) {
  return (
    <>
      <Box borderWidth="1px" borderRadius="lg" overflow="hidden">
        <Box>
          <Center>
            <Img
              h="175px"
              src={props.cdn_url}
              fallbackSrc={props.spaces_url}
              alt={props.file_name}
            />
          </Center>
        </Box>

        <Divider />

        <Flex>
          <Box p="2">
            <Skeleton isLoaded={!props.loading}>
              <Heading size="md">{props.file_name}</Heading>
            </Skeleton>
            <Text as={SkeletonText} fontSize="xs" isLoaded={!props.loading}>
              {formatDate(props.last_modified)} &bull;{' '}
              {formatFileSize(props.size)}
            </Text>
          </Box>
          <Spacer />
          <Box>
            <Menu>
              <MenuButton
                as={Button}
                rightIcon={<BsChevronDown />}
                m="2"
                isLoading={props.loading}
                size="xs"
              >
                Actions
              </MenuButton>
              <MenuList>
                <MenuItem>
                  <RouterLink
                    to={{ pathname: props.cdn_url }}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    Open in New Tab
                  </RouterLink>
                </MenuItem>
                <MenuItem
                  as={RouterLink}
                  to={{ pathname: `${props.cdn_url}?download=true` }}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Download
                </MenuItem>
                <MenuItem>Delete</MenuItem>
              </MenuList>
            </Menu>
          </Box>
        </Flex>
      </Box>
    </>
  );
}
