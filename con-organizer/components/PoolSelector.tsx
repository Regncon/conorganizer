'use client';

import { useState } from 'react';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';
import { Pool } from '@/lib/enums';
import { a11yProps } from '@/lib/helpers';

type Props = {
    handlePoolChange: (pool: Pool) => void;
};

const PoolSelector = ({ handlePoolChange }: Props) => {
    const [value, setValue] = useState<Pool>(Pool.FridayEvening);

    const handleChange = (event: React.SyntheticEvent, newValue: Pool) => {
        switch (newValue) {
            case Pool.FridayEvening:
                setValue(Pool.FridayEvening);
                handlePoolChange(Pool.FridayEvening);
                break;
            case Pool.SaturdayMorning:
                setValue(Pool.SaturdayMorning);
                handlePoolChange(Pool.SaturdayMorning);
                break;
            case Pool.SaturdayEvening:
                setValue(Pool.SaturdayEvening);
                handlePoolChange(Pool.SaturdayEvening);
                break;
            case Pool.SundayMorning:
                setValue(Pool.SundayMorning);
                handlePoolChange(Pool.SundayMorning);
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
