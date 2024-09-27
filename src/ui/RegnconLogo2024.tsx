'use client';
import React from 'react';
import { Box } from '@mui/material';

type Props = {
    spin?: boolean;
    size?: 'small' | 'medium' | 'large';
};

const RegnconLogo2024 = ({ spin = false, size = 'medium' }: Props) => {
    const getSize = () => {
        switch (size) {
            case 'small':
                return 42;
            case 'large':
                return 400;
            default:
                return 100;
        }
    };

    return (
        <Box
            component="img"
            src="/RegnCon2024Logo1_2_2.svg"
            alt="Regncon 2024 Logo"
            sx={{
                width: '90vw',
                maxWidth: getSize(),
                height: '90vw',
                maxHeight: getSize(),
                animation: spin ? 'spin 240s linear infinite' : 'none',
                '@keyframes spin': {
                    from: {
                        transform: 'rotate(0deg)',
                    },
                    to: {
                        transform: 'rotate(-360deg)',
                    },
                },
            }}
        />
    );
};

export default RegnconLogo2024;
