'use client';

import { ErrorBoundary } from 'react-error-boundary';
import EventBoundary from '@/components/ErrorBoundaries/EventBoundary';
import MainNavigator from '@/components/MainNavigator';
import Event from './Event';
type Props = {
    params: { id: string };
};
const page = ({ params }: Props) => {
    const { id } = params;
    return (
        <>
            <ErrorBoundary FallbackComponent={EventBoundary}>
                <Event id={id} />
            </ErrorBoundary>
            <MainNavigator />
        </>
    );
};
export default page;
