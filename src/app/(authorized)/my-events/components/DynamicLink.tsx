'use client';
import Box from '@mui/material/Box';
import type { Route } from 'next';
import { useRouter } from 'next/navigation';

import { useEffect, type PropsWithChildren } from 'react';

type Props = { docId: string; disable?: boolean } & PropsWithChildren;

const DynamicLink = ({ children, docId, disable }: Props) => {
    const router = useRouter();
    const createEventHref = `/event/create/${docId}` as Route;
    useEffect(() => {
        if (disable) {
            return;
        }
        router.prefetch(createEventHref as Route);
    }, []);
    return (
        <Box
            onClick={() => {
                if (disable) {
                    return;
                }
                router.push(createEventHref);
            }}
        >
            {children}
        </Box>
    );
};

export default DynamicLink;
