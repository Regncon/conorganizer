'use client';

import { migrateParticipantAndInterest } from '$lib/serverActions/Migration';
import { useEffect } from 'react';

export default function Test() {
    useEffect(() => {
        migrateParticipantAndInterest().then((res) => {
            console.log(res);
        });
    });
    return null;
}
