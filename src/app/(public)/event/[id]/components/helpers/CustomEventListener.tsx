'use client';

import { useEffect } from 'react';

const CustomEventListener = () => {
    useEffect(() => {
        const myParticipantsChangedFalse = new CustomEvent('my-participants-changed', {
            detail: { loading: false },
        });
        document.dispatchEvent(myParticipantsChangedFalse);
    });

    return null;
};
export default CustomEventListener;
