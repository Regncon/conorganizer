'use client';
import Box from '@mui/material/Box';
import type { Route } from 'next';
import { useRouter } from 'next/navigation';

import { useEffect, type PropsWithChildren } from 'react';

type Props = { docId: string } & PropsWithChildren;

const DynamicLink = ({ children, docId }: Props) => {
    const router = useRouter();
    const createEventHref = `/event/create/${docId}` as Route;
    useEffect(() => {
        router.prefetch(createEventHref as Route);
    }, []);
    return (
        <Box
            onClick={() => {
                router.push(createEventHref);
            }}
        >
            {children}
        </Box>
    );
};

export default DynamicLink;
