import React, { useEffect, useState } from 'react';
import { Helmet } from 'react-helmet';
import ContainerCentered from '../components/ContainerCentered';
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
import { Link as RouterLink, useHistory } from 'react-router-dom';
import { GoSignOut } from 'react-icons/go';
import axios from 'axios';
import { useAuth } from '../contexts/AuthProvider';
import { useToken } from '../contexts/TokenContext';

export default function Profile() {
  const auth = useAuth();
  const history = useHistory();

  const [showToken, setShowToken] = useState(false);
  const tokenState = useToken(auth.user!.uid);
  const [token, setToken] = useState(tokenState.token);
  const { hasCopied, onCopy } = useClipboard(token);

  const handleShowToken = () => setShowToken((show) => !show);

  const handleClick = () => {
    auth.logout!(history).catch((error) => {
      console.log(error.message);
    });
  };

  const handleTokenGeneration = async () => {
    try {
      const apiURL = import.meta.env.SNOWPACK_PUBLIC_API_URL;
      const token = await auth.user!.getIdToken(true);
      await axios.post(`${apiURL}/token`, { token });
    } catch (error) {
      console.log(error.message);
    }
  };

  useEffect(() => {
    setToken(tokenState.token);
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
                  value={token}
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
                  <Button h="1.75rem" size="sm" onClick={handleShowToken}>
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