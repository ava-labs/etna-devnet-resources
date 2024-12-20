import PasteKeys from "./PasteKeys";
import GenerateKeys from "./GenerateKeys";
import Genesis from "./Genesis";
import Steps from "./Steps";
import { stepList, useWizardStore } from "./store";
import CreateL1 from "./CreateL1";
import { TechInfo } from "./TechInfo";
import LaunchValidators from "./LaunchValidators";


const stepComponents: Record<keyof typeof stepList, React.ReactNode> = {
    'genesis': <Genesis />,
    'generate-keys': <GenerateKeys />,
    'paste-keys': <PasteKeys />,
    'create-l1': <CreateL1 />,
    "launch-validators": <LaunchValidators />,
}


export default function L1Wizard() {
    const { currentStep } = useWizardStore()

    return (
        <>
            <div className="flex container mx-auto max-w-6xl py-8">
                <div className="w-80 p-4 shrink-0">
                    <Steps />
                </div>
                <div className="flex-1 pl-4 min-w-0">
                    <div className="h-full">
                        {stepComponents[currentStep]}
                    </div>
                </div>

            </div>
            <div className=" p-4 text-center text-xs">
                <div className="text-gray-600">
                    <TechInfo />
                </div>
            </div>
        </>
    );
}
