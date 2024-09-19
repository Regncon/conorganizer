'use client';

import MuiMarkdown from 'mui-markdown';

type Props = {
    description: string;
};

const MuiMarkdownClient = ({ description }: Props) => {
    return <MuiMarkdown>{description || '# Lang beskrivelse'}</MuiMarkdown>;
};

export default MuiMarkdownClient;
