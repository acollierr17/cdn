import React from 'react';
import { Link as ChakraLink } from '@chakra-ui/react';
import { Link as RouterLink } from 'react-router-dom';

export default function Link({ children, ...rest }: any) {
  return (
    <ChakraLink as={RouterLink} {...rest}>
      {children}
    </ChakraLink>
  );
}
