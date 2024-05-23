type ErrorBoundaryProps = ({ error, reset }: { error: Error & { digest?: string }; reset: () => void }) => JSX.Element;
