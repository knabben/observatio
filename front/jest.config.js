import nextJest from 'next/jest.js';

const createJestConfig = nextJest();


const config = {
  testEnvironment: 'jest-environment-jsdom',
  setupFilesAfterEnv: ['<rootDir>/jest.setup.tsx'],
  preset: 'ts-jest',
  moduleNameMapper: {
    '^@/app/(.*)$': '<rootDir>/app/$1',
    '^@/fonts$': '<rootDir>/app/styles/fonts',
  },
};

export default createJestConfig(config);
