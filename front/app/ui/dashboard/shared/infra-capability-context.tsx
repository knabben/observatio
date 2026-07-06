'use client';

import React, {createContext, useContext} from 'react';
import {InfrastructureCapability, emptyInfrastructureCapability} from '@/app/lib/data';

/**
 * Makes the environment's detected infrastructure capability (fetched once by the
 * Clusters/Machines tab shells) available to any nested row/badge without prop-drilling
 * or re-fetching per row.
 */
const InfraCapabilityContext = createContext<InfrastructureCapability>(emptyInfrastructureCapability);

export const InfraCapabilityProvider = InfraCapabilityContext.Provider;

export function useInfraCapability(): InfrastructureCapability {
  return useContext(InfraCapabilityContext);
}
