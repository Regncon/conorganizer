'use client';

import { useMediaQuery } from '@mui/material';
import type { ComponentPropsWithoutRef, PropsWithChildren } from 'react';

type Props = {};

const BigMediaQueryWrapper = ({ children }: PropsWithChildren<Props>) => {
    const isBigScreen = useMediaQuery('(min-width:600px)');
    return isBigScreen ? children : null;
};

export default BigMediaQueryWrapper;
