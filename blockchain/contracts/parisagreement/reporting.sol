pragma solidity ^0.5.10;

    /**
    Contract for to report intermediate accumulated GHG emissions from the countries Paris Agreement participants.
    Accumulated results from databases will commit automatically from Openclimate backend.
    */
contract Reporting {
    /**
    Struct to keep reporting entries
    */
    struct country_report {
        bytes32 country_name;
        // GHG mitigation goals
        int CO2; // (required) in metric tonnes
        int CH4; // (required) in metric tonnes
        int N2O; // (required) in metric tonnes
        int HFCs; // in metric tonnes
        int PFCs; // in metric tonnes
        int SF6; // in metric tonnes
        int NF3; // in metric tonnes
        int AltEnergy; // alternative/renewable energy usage in MWh
        // Finance
        int Mobilization; // (required) in billions USD.
        int ContribGreenFund; // (required) Contribution to Green Climate Fund in mibillionsllions USD
        int BilateralLoan; //  Bilateral loan to developing country Party billions USD
        // util variables
        uint timeStamp; // timestamp of the latest report
    }

    /**
    Mapping to store reports by address.
    */
    mapping(address=> mapping(uint => country_report)) private Reports;
    mapping(address=> uint[]) private timeStamps;
    
    /**
       Function to insert a new report. It increments values.
    */
    function insert_report(
        bytes32 country_name,
        int CO2,
        int CH4,
        int N2O,
        int HFCs,
        int PFCs,
        int SF6,
        int NF3,
        int AltEnergy,
        int Mobilization,
        int ContribGreenFund,
        int BilateralLoan,
        uint timeStamp)
    public
    {
        // sanity check if there is timestamp already exist
        if(Reports[msg.sender][timeStamp].timeStamp==timeStamp) revert();
        // Next timestamp should be higher than the prevoius
        if(getLastStamp(msg.sender)>=timeStamp) revert();
        // Push new timeStamp
        timeStamps[msg.sender].push(timeStamp);
        uint lastTimestamp = getLastStamp(msg.sender);
        // insert values
        Reports[msg.sender][timeStamp].country_name      = country_name;
        Reports[msg.sender][timeStamp].CO2               = Reports[msg.sender][lastTimestamp].CO2+CO2;
        Reports[msg.sender][timeStamp].CH4               = Reports[msg.sender][lastTimestamp].CH4+CH4;
        Reports[msg.sender][timeStamp].N2O               = Reports[msg.sender][lastTimestamp].N2O+N2O;
        Reports[msg.sender][timeStamp].HFCs              = Reports[msg.sender][lastTimestamp].HFCs+HFCs;
        Reports[msg.sender][timeStamp].PFCs              = Reports[msg.sender][lastTimestamp].PFCs+PFCs;
        Reports[msg.sender][timeStamp].SF6               = Reports[msg.sender][lastTimestamp].SF6+SF6;
        Reports[msg.sender][timeStamp].NF3               = Reports[msg.sender][lastTimestamp].NF3+NF3;
        Reports[msg.sender][timeStamp].AltEnergy         = Reports[msg.sender][lastTimestamp].AltEnergy+AltEnergy;
        Reports[msg.sender][timeStamp].Mobilization      = Reports[msg.sender][lastTimestamp].Mobilization+Mobilization;
        Reports[msg.sender][timeStamp].ContribGreenFund  = Reports[msg.sender][lastTimestamp].ContribGreenFund+ContribGreenFund;
        Reports[msg.sender][timeStamp].BilateralLoan     = Reports[msg.sender][lastTimestamp].BilateralLoan+BilateralLoan;
        Reports[msg.sender][timeStamp].timeStamp         = timeStamp;
    }
    
    function getAllTimeStamps () view public returns (uint[] memory timestamps) {
        return timeStamps[msg.sender];
    }
    
    function getLastStamp (address countryAddr) view public returns (uint lastTimestamp){
        uint[] memory list = timeStamps[countryAddr];
        uint length = timeStamps[countryAddr].length;
        if (length == 0) return 0;
        else return list[length-1];
    }
    
    function getPrevoiusStamp (address countryAddr) view public returns (uint prevousTimestamp){
        uint[] memory list = timeStamps[countryAddr];
        uint length = timeStamps[countryAddr].length;
        if (length == 1) return 0;
        else return list[length-2];
    }
    
    function getLastReport(address countryAddr) 
    view 
    public 
    returns 
    (
        bytes32 country_name,
        int CO2,
        int CH4,
        int N2O,
        int HFCs,
        int PFCs,
        int SF6,
        int NF3,
        int AltEnergy,
        int Mobilization,
        int ContribGreenFund,
        int BilateralLoan,
        uint timeStamp
        ) 
    
    {
        uint  lastTimestamp = getLastStamp(countryAddr);
        return (
            Reports[msg.sender][lastTimestamp].country_name,
            Reports[msg.sender][lastTimestamp].CO2,
            Reports[msg.sender][lastTimestamp].CH4,
            Reports[msg.sender][lastTimestamp].N2O,
            Reports[msg.sender][lastTimestamp].HFCs,
            Reports[msg.sender][lastTimestamp].PFCs,
            Reports[msg.sender][lastTimestamp].SF6,
            Reports[msg.sender][lastTimestamp].NF3,
            Reports[msg.sender][lastTimestamp].AltEnergy,
            Reports[msg.sender][lastTimestamp].Mobilization,
            Reports[msg.sender][lastTimestamp].ContribGreenFund,
            Reports[msg.sender][lastTimestamp].BilateralLoan,
            Reports[msg.sender][lastTimestamp].timeStamp
            );
    }
    
    function getLastCO2(address countryAddr) 
    view 
    public 
    returns 
    (
        int CO2,
        uint timeStamp
        ) 
    
    {
        uint  lastTimestamp = getLastStamp(countryAddr);
        return (
            Reports[countryAddr][lastTimestamp].CO2,
            Reports[countryAddr][lastTimestamp].timeStamp
            );
    }
    
    function getIncrementCO2(address countryAddr) 
    view 
    public 
    returns 
    (
        int CO2
        ) 
    {
        uint  lastTimestamp = getLastStamp(countryAddr);
        uint  prevousTimestamp = getPrevoiusStamp(countryAddr);
        return (
            Reports[countryAddr][lastTimestamp].CO2-Reports[countryAddr][prevousTimestamp].CO2
            );
    }
}
