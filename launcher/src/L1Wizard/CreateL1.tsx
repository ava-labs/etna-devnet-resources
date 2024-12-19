import NextPrev from "./ui/NextPrev";

export default function CreateL1() {
    return (
        <div>
            <h1>Create an L1</h1>
            <NextPrev nextDisabled={false} currentStepName="create-l1" />
        </div>
    )
}
