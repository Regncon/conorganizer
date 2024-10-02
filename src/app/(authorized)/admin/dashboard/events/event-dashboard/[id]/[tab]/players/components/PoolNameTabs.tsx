'use client';

import { PoolName } from '$lib/enums';
import { Tabs, Tab, debounce } from '@mui/material';
import type { Route } from 'next';
import { useRouter, useSearchParams } from 'next/navigation';
import { useCallback, useEffect, useState, type PropsWithChildren } from 'react';

function a11yProps(index: PoolName) {
    return {
        id: `PoolName-tab-${index}`,
        'aria-controls': `PoolName-tab-${index}`,
    };
}

type Props = {
    id: string;
};

const PoolNameTabs = ({ id }: Props) => {
    console.log('tab', id);

    const router = useRouter();

    const searchParams = useSearchParams();
    const tab = (searchParams.get('active-tab') ?? 'fridayEvening') as PoolName;

    const [value, setValue] = useState<PoolName>(tab);

    const handleChange = (event: React.SyntheticEvent, activeTab: PoolName) => {
        setValue(activeTab);
        router.replace(('?active-tab=' + activeTab) as Route);
    };

    useEffect(() => {
        router.prefetch(`/admin/dashboard/events/event-dashboard/${id}/players?active-tab=${PoolName.fridayEvening}`);
        router.prefetch(`/admin/dashboard/events/event-dashboard/${id}/players?active-tab=${PoolName.saturdayMorning}`);
        router.prefetch(`/admin/dashboard/events/event-dashboard/${id}/players?active-tab=${PoolName.saturdayEvening}`);
        router.prefetch(`/admin/dashboard/events/event-dashboard/${id}/players?active-tab=${PoolName.sundayMorning}`);
    }, []);

    return (
        <Tabs value={value} onChange={handleChange} aria-label="PoolName selector">
            <Tab label="Fredag Kveld" {...a11yProps(PoolName.fridayEvening)} value={PoolName.fridayEvening} />
            <Tab label="Lørdag Morgen" {...a11yProps(PoolName.saturdayMorning)} value={PoolName.saturdayMorning} />
            <Tab label="Lørdag Kveld" {...a11yProps(PoolName.saturdayEvening)} value={PoolName.saturdayEvening} />
            <Tab label=" Søndag Morgen" {...a11yProps(PoolName.sundayMorning)} value={PoolName.sundayMorning} />
        </Tabs>
    );
};

export default PoolNameTabs;
