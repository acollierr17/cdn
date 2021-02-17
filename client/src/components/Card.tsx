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
import { BsChevronDown, BsChevronUp } from 'react-icons/bs';
import type { ImageResult } from '../api';
import { formatDate, formatFileSize } from '../functions';
import Link from '../components/Link';
import DeleteImage from './DeleteImage';

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
            <Skeleton isLoaded={!props.loading}>
              <Img
                h="175px"
                maxW="inherit"
                src={props.cdn_url}
                fallbackSrc={props.spaces_url}
                alt={props.file_name}
              />
            </Skeleton>
          </Center>
        </Box>

        <Divider />

        <Flex>
          <Box p="2">
            <Skeleton isLoaded={!props.loading}>
              <Heading
                as={Link}
                size="md"
                to={{ pathname: props.cdn_url }}
                target="_blank"
                rel="noopener noreferrer"
              >
                {props.file_name}
              </Heading>
            </Skeleton>
            <Text
              as={SkeletonText}
              fontSize="xs"
              isLoaded={!props.loading}
              h="40px"
              w="110%"
            >
              {formatDate(props.last_modified)} &bull;{' '}
              {formatFileSize(props.size)}
            </Text>
          </Box>
          <Spacer />
          <Box>
            <Menu>
              {({ isOpen }) => (
                <>
                  <MenuButton
                    as={Button}
                    rightIcon={isOpen ? <BsChevronUp /> : <BsChevronDown />}
                    m="2"
                    isLoading={props.loading}
                    size="sm"
                  >
                    {isOpen ? 'Close' : 'Open'}
                  </MenuButton>
                  <MenuList>
                    <MenuItem
                      as={RouterLink}
                      to={{ pathname: `${props.cdn_url}?download=true` }}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      Download
                    </MenuItem>
                    <DeleteImage
                      fileName={props.file_name}
                      token={props.token}
                    />
                  </MenuList>
                </>
              )}
            </Menu>
          </Box>
        </Flex>
      </Box>
    </>
  );
}
