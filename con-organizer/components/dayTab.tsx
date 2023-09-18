'use client';

import { useState } from 'react';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import { useTheme } from '@mui/material/styles';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';
import Typography from '@mui/material/Typography';
import { Pool } from '@/lib/enums';

interface TabPanelProps {
    children?: React.ReactNode;
    dir?: string;
    index: number;
    value: number;
}

function TabPanel(props: TabPanelProps) {
    const { children, value, index, ...other } = props;

    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`full-width-tabpanel-${index}`}
            aria-labelledby={`full-width-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box sx={{ p: 3 }}>
                    <Typography>{children}</Typography>
                </Box>
            )}
        </div>
    );
}

function a11yProps(index: number) {
    return {
        id: `full-width-tab-${index}`,
        'aria-controls': `full-width-tabpanel-${index}`,
    };
}

type Props = {
    handlePoolChange?: (pool: Pool) => void;
};

export default function FullWidthTabs({ handlePoolChange }: Props) {
    const theme = useTheme();
    const [value, setValue] = useState<Pool>(Pool.FridayEvening);

    const handleChange = (event: React.SyntheticEvent, newValue: Pool) => {
        switch (newValue) {
            case Pool.FridayEvening:
                setValue(Pool.FridayEvening);
                break;
            case Pool.SaturdayMorning:
                setValue(Pool.SaturdayMorning);
                break;
            case Pool.SaturdayEvening:
                setValue(Pool.SaturdayEvening);
                break;
            case Pool.SundayMorning:
                setValue(Pool.SundayMorning);
                break;
            default:
                setValue(Pool.none);
                break;
        }
    };

    return (
        <Box sx={{ bgcolor: 'background.paper' }}>
            <AppBar position="fixed">
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
}
