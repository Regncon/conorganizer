'use client';

import { useCallback, useEffect, useState } from 'react';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';
import { Route } from 'next';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { a11yProps } from '@/lib/helpers';
import { Pool } from '@/models/enums';


type Props = {
    handlePoolChange: (pool: Pool) => void;
};

const PoolSelector = ({ handlePoolChange }: Props) => {
    const [value, setValue] = useState<Pool>(Pool.FridayEvening);

    const router = useRouter();
    const pathname = usePathname();
    const searchParams = useSearchParams();

    const createQueryString = useCallback(
        (name: string, value: string) => {
            const params = new URLSearchParams(searchParams);
            params.set(name, value);

            return params.toString();
        },
        [searchParams]
    );

    const search = searchParams.get('pool') as keyof typeof Pool;
    useEffect(() => {
        if (search) {
            setValue(Pool[search]);
        }
    }
    , [search]);


    const handleChange = (event: React.SyntheticEvent, newValue: Pool) => {
        switch (newValue) {
            case Pool.FridayEvening:
                setValue(Pool.FridayEvening);
                handlePoolChange(Pool.FridayEvening);
                router.push(`${pathname}?${createQueryString('pool', 'FridayEvening')}` as Route);
                break;
            case Pool.SaturdayMorning:
                setValue(Pool.SaturdayMorning);
                handlePoolChange(Pool.SaturdayMorning);
                router.push(`${pathname}?${createQueryString('pool', 'SaturdayMorning')}` as Route);
                break;
            case Pool.SaturdayEvening:
                setValue(Pool.SaturdayEvening);
                handlePoolChange(Pool.SaturdayEvening);
                router.push(`${pathname}?${createQueryString('pool', 'SaturdayEvening')}` as Route);
                break;
            case Pool.SundayMorning:
                setValue(Pool.SundayMorning);
                handlePoolChange(Pool.SundayMorning);
                router.push(`${pathname}?${createQueryString('pool', 'SundayMorning')}` as Route);
                break;
            default:
                setValue(Pool.none);
                break;
        }
        
    };

    return (
        <Box sx={{ bgcolor: 'background.paper' }}>
            <AppBar position="sticky" sx={{ width: '100vw' }}>
                <Tabs
                    value={value}
                    onChange={handleChange}
                    indicatorColor="secondary"
                    textColor="inherit"
                    variant="fullWidth"
                    aria-label="full width tabs example"
                >
                    <Tab value={Pool.FridayEvening} label={Pool.FridayEvening} {...a11yProps(0)} />
                    <Tab value={Pool.SaturdayMorning} label={Pool.SaturdayMorning} {...a11yProps(1)} />
                    <Tab value={Pool.SaturdayEvening} label={Pool.SaturdayEvening} {...a11yProps(2)} />
                    <Tab value={Pool.SundayMorning} label={Pool.SundayMorning} {...a11yProps(2)} />
                </Tabs>
            </AppBar>
        </Box>
    );
};

export default PoolSelector;
