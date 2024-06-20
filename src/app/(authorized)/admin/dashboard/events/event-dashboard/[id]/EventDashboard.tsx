'use client';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Box from '@mui/material/Box';
import { useState } from 'react';
import Settings from './Settings';
import Edit from './Edit';

interface TabPanelProps {
    children?: React.ReactNode;
    index: number;
    value: number;
}

function CustomTabPanel(props: TabPanelProps) {
    const { children, value, index, ...other } = props;

    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`simple-tabpanel-${index}`}
            aria-labelledby={`simple-tab-${index}`}
            {...other}
        >
            {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
        </div>
    );
}

function a11yProps(index: number) {
    return {
        id: `simple-tab-${index}`,
        'aria-controls': `simple-tabpanel-${index}`,
    };
}

export default function EventDashboard() {
    const [value, setValue] = useState(0);

    const handleChange = (event: React.SyntheticEvent, newValue: number) => {
        setValue(newValue);
    };

    return (
        <Box sx={{ width: '100%' }}>
            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                <Tabs value={value} onChange={handleChange} aria-label="basic tabs example">
                    <Tab label="Spillere" {...a11yProps(0)} />
                    <Tab label="Innstillinger" {...a11yProps(1)} />
                    <Tab label="Rom" {...a11yProps(2)} />
                    <Tab label="Rediger" {...a11yProps(3)} />
                </Tabs>
            </Box>
            <CustomTabPanel value={value} index={0}>
                Item 0
            </CustomTabPanel>
            <CustomTabPanel value={value} index={1}>
                <Settings />
            </CustomTabPanel>
            <CustomTabPanel value={value} index={2}>
                Item 10
            </CustomTabPanel>
            <CustomTabPanel value={value} index={3}>
                <Edit />
            </CustomTabPanel>
        </Box>
    );
}
