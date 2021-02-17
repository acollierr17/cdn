import React, { useEffect, useState, useRef } from 'react';
import { Helmet } from 'react-helmet';
import {
  FormControl,
  FormLabel,
  Heading,
  Input,
  InputGroup,
  InputRightElement,
  Box,
  Divider,
  Button,
  Text,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  Flex,
  Spacer,
  useClipboard,
} from '@chakra-ui/react';
import ContainerCentered from '../components/ContainerCentered';
import { Link as RouterLink, useHistory } from 'react-router-dom';
import { GoSignOut } from 'react-icons/go';
import { useAuth } from '../contexts/AuthContext';
import { useToken } from '../contexts/TokenContext';
import { generateToken } from '../api';

export default function Profile() {
  const auth = useAuth();
  const history = useHistory();

  const [showToken, setShowToken] = useState(false);
  const tokenState = useToken(auth.user!.uid);
  const [token, setToken] = useState(tokenState.token);
  const { hasCopied, onCopy } = useClipboard(token);
  const tokenRef = useRef<any>();

  const toggleTokenVisibility = () => setShowToken((show) => !show);

  const handleClick = () => {
    auth.logout!(history).catch((error) => {
      console.log(error.message);
    });
  };

  const handleTokenGeneration = async () => {
    try {
      await generateToken(auth);
    } catch (error) {
      console.log(error.message);
    }
  };

  useEffect(() => {
    setToken(tokenState.token);
    tokenRef.current.value = token;
  }, [tokenState]);

  return (
    <>
      <Helmet>
        <title>Profile</title>
        <meta property="og:title" content="Profile" />
      </Helmet>
      <ContainerCentered>
        <Flex>
          <Breadcrumb>
            <BreadcrumbItem as={Heading} pb={1}>
              <BreadcrumbLink as={RouterLink} to="/">
                Dashboard
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbItem as={Heading} pb={1} isCurrentPage>
              <BreadcrumbLink as={RouterLink} to="/profile">
                Profile
              </BreadcrumbLink>
            </BreadcrumbItem>
          </Breadcrumb>
          <Spacer />
          <Button
            leftIcon={<GoSignOut />}
            colorScheme="orange"
            onClick={handleClick}
          >
            Sign Out
          </Button>
        </Flex>
        <Divider />
        <Box pt={4}>
          <form>
            <FormControl id="email">
              <FormLabel>Email address</FormLabel>
              <Input
                type="email"
                name="email"
                value={auth.user!.email!}
                isDisabled
              />
            </FormControl>

            <FormControl id="access-token" pt={4}>
              <FormLabel>Access Token</FormLabel>
              <InputGroup size="md">
                <Text
                  as={Input}
                  pr="13.5rem"
                  type={showToken ? 'text' : 'password'}
                  ref={tokenRef}
                  isDisabled
                />
                <InputRightElement width="5rem" mr="126px">
                  <Button h="1.75rem" size="sm" onClick={handleTokenGeneration}>
                    Reset
                  </Button>
                </InputRightElement>
                <InputRightElement width="4.5rem" mr="65px">
                  <Button h="1.75rem" size="sm" onClick={onCopy}>
                    {hasCopied ? 'Copied' : 'Copy'}
                  </Button>
                </InputRightElement>
                <InputRightElement width="4.5rem">
                  <Button h="1.75rem" size="sm" onClick={toggleTokenVisibility}>
                    {showToken ? 'Hide' : 'Show'}
                  </Button>
                </InputRightElement>
              </InputGroup>
            </FormControl>
          </form>
        </Box>
      </ContainerCentered>
    </>
  );
}
