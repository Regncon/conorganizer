"use client";

import * as React from 'react';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import { useTheme } from '@mui/material/styles';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';
import Typography from '@mui/material/Typography';
import { pool } from '@/lib/enums';

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
  handlePoolChange?: (pool: pool) => void;
};

export default function FullWidthTabs({ handlePoolChange }: Props) {
  const theme = useTheme();
  const [value, setValue] = React.useState<pool>(pool.FirdayEvening);

  const handleChange = (event: React.SyntheticEvent, newValue: number) => {
    if (newValue === 0) {
      setValue(pool.FirdayEvening);
    }
    if (newValue === 1) {
      setValue(pool.SaturdayMorning);
    }
    if (newValue === 2) {
      setValue(pool.SaturdayEvening);
    }
    if (newValue === 3) {
      setValue(pool.SundayMorning);
    }
  };

  return (
    <Box sx={{ bgcolor: 'background.paper' }}>
      <AppBar position="static">
        <Tabs
          value={value}
          onChange={handleChange}
          indicatorColor="secondary"
          textColor="inherit"
          variant="fullWidth"
          aria-label="full width tabs example"
        >
          <Tab label="Fredag kveld" {...a11yProps(0)} />
          <Tab label="Lørdag morgen" {...a11yProps(1)} />
          <Tab label="Lørdag kveld" {...a11yProps(2)} />
          <Tab label="Søndag morgen" {...a11yProps(2)} />
        </Tabs>
      </AppBar>
    </Box>
  );
}