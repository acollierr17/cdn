import React from 'react';
import { Code } from '@chakra-ui/react';

interface ConfigProps {
  config: string;
}

export default function ShareXConfig({ config }: ConfigProps) {
  return (
    <Code
      display="block"
      whiteSpace="pre"
      overflow="scroll"
      children={config}
    />
  );
}
