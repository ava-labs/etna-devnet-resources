import CreateL1 from "./CreateL1";
import GenerateKeys from "./GenerateKeys";
import Genesis from "./Genesis";
import Steps from "./Steps";
import { stepList, useWizardStore } from "./store";


const stepComponents: Record<keyof typeof stepList, React.ReactNode> = {
    'genesis': <Genesis />,
    'generate-keys': <GenerateKeys />,
    'create-l1': <CreateL1 />,
}


export default function L1Wizard() {
    const { currentStep } = useWizardStore()

    return (
        <div className="flex">
            <div className="w-80 p-4">
                <Steps />
            </div>
            <div className="flex-1 pl-4">
                <div className="h-full">
                    {stepComponents[currentStep]}
                </div>
            </div>
        </div>
    );
}
