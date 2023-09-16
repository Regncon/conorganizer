"use client"

type Props = {
    error: string;
}

const error = ({ error }: Props) => {
    return (
        <div>
            <h1>
                Error: {error}
            </h1>
        </div>
    );
}

export default error;