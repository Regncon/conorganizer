import CardBase from '$app/(authorized)/dashboard/CardBase';

const FormSubmissions = async () => {
    return (
        <CardBase
            href="/admin/dashboard/form-submissions"
            subTitle="Trykk for å gå til innsendte skjemaer"
            img="/form-submissions-thumbnail.webp"
            imgAlt="Innsendte skjemaer"
            title="Innsendte skjemaer"
        />
    );
};

export default FormSubmissions;
