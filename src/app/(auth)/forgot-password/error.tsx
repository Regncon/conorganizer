'use client';
import LoginErrorBoundary from '$ui/ErrorBoundary/LoginErrorBoundary';

const Error: ErrorBoundaryProps = ({ error, reset }) => {
    return <LoginErrorBoundary error={error} reset={reset} />;
};

export default Error;
