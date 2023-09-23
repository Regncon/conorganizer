'use client';

import MainNavigator from '@/components/MainNavigator';
import Event from './Event';
type Props = {
    params: { id: string };
};
const page = ({ params }: Props) => {
    const { id } = params;
    return (
        <>
            <Event id={id} />
            <MainNavigator />
        </>
    );
};
export default page;
