'use client';

import { useMediaQuery } from '@mui/material';
import type { ComponentPropsWithoutRef, PropsWithChildren } from 'react';

type Props = {};

const SmallMediaQueryWrapper = ({ children }: PropsWithChildren<Props>) => {
    const isBigScreen = useMediaQuery('(max-width:599px)');
    return isBigScreen ? children : null;
};

export default SmallMediaQueryWrapper;
