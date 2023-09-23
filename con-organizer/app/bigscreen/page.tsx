'use client';
import { ErrorBoundary } from 'react-error-boundary';
import BigScreenBoundary from '@/components/ErrorBoundaries/BigScreenBoundary';
import BigScreen from './bigScreen';

const page = () => {
    return (
        <ErrorBoundary FallbackComponent={BigScreenBoundary}>
            <BigScreen />
        </ErrorBoundary>
    );
};

export default page;
