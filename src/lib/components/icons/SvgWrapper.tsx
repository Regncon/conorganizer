'use client';
import type { ColorProp } from '$app/(public)/components/lib/helpers/icons';
import { SvgIcon, useTheme, type Palette } from '@mui/material';
import type { PropsWithChildren } from 'react';

export type Props = {
    color?: ColorProp;
    size?: 'small' | 'medium' | 'large' | 'inherit';
    chipMargin?: boolean;
};

const SvgWrapper = ({ children, color = 'primary', size = 'small', chipMargin = true }: PropsWithChildren<Props>) => {
    const theme = useTheme();
    return (
        <SvgIcon
            sx={{
                marginInlineStart: chipMargin ? '0.7rem' : '0',
                fontSize:
                    size === 'small' ? '1.7rem'
                    : size === 'medium' ? '2rem'
                    : size === 'large' ? '3rem'
                    : 'inherit',
                color: 'inherit',
                fill: theme.palette[color].main,
            }}
        >
            {children}
        </SvgIcon>
    );
};

export default SvgWrapper;
