'use client';

import { useMediaQuery } from '@mui/material';
import type { ComponentPropsWithoutRef, PropsWithChildren } from 'react';

type Props = {};

const SmallMediaQueryWrapper = ({ children }: PropsWithChildren<Props>) => {
    const isBigScreen = useMediaQuery('(max-width:676px)');
    return isBigScreen ? children : null;
};

export default SmallMediaQueryWrapper;
