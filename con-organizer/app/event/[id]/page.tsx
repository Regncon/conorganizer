import Event from './Event';
type Props = {
    params: { id: string };
};

const page = ({ params }: Props) => {
    const { id } = params;
    console.log(id);

    return <Event id={id} />;
};

export default page;
