pragma solidity ^0.5.10;

contract DetermedContributions {
    /**
    According to Paris Agreement (PA) each country shall prepare, communicate and maintain successive nationally determined
    contributions (NDC) that it intends to achieve (GHG emissions). The struct allows to keep track of NDC for every country.

    Katowice Climate define following guidelines for reporting:
     - national inventory report on anthropogenic emissions by sources and removals by sinks of greenhouse gases; (first biennial report under new MPGs due December 2024)
        1. 2006 IPCC guidelines
        2. 100 year gwp CO2e
        3. seven gases (CO2, CH4, N2O, HFCs, PFCs, SF6 and NF3);
    */
    struct country_NDC {
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
        // Adaptation goals
        // Capacity building
        // Finance
        int Mobilization; // (required) in billions USD.
        int ContribGreenFund; // (required) Contribution to Green Climate Fund in mibillionsllions USD
        int BilateralLoan; //  Bilateral loan to developing country Party billions USD
        // Technology transfer
        // Transparency
        // util variables
        uint timeTarget; // set a time target
        uint index;
    }


    /**
    Creating mapping to store countries goals. This allows to save/delete/update pledges for each country.
    */
    mapping(address => country_NDC) private NDCs;
    address[] private countryIndex;


    // /**
    // Util function to emit variable after inserting/updating/deleting
    // */
    // event LogNewNDC   (address indexed country_address, uint index, bytes32 country_name, int CO2, int CH4, int N2O, int HFCs, int PFCs, int SF6, int NF3, int AltEnergy, int Mobilization, int ContribGreenFund, int BilateralLoan, uint timeTarget);
    // event LogUpdateNDC(address indexed country_address, uint index, bytes32 country_name, int CO2, int CH4, int N2O, int HFCs, int PFCs, int SF6, int NF3, int AltEnergy, int Mobilization, int ContribGreenFund, int BilateralLoan, uint timeTarget);
    // event LogDeleteNDC(address indexed country_address, uint index);

    /**
    Util function toto check if the public key of the country in the list
    */
    function isCounrty(address counrty)
    public
    view
    returns(bool isIndeed)
    {
        if(countryIndex.length == 0) return false;
        return (countryIndex[NDCs[counrty].index] == counrty);
    }

    /**
       Function to insert a new NDC from a country address
    */
    function insertNDC(
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
        uint timeTarget)
    public
    returns(uint index)
    {
        if(isCounrty(msg.sender)) revert();
        NDCs[msg.sender].country_name      = country_name;
        NDCs[msg.sender].CO2               = CO2;
        NDCs[msg.sender].CH4               = CH4;
        NDCs[msg.sender].N2O               = N2O;
        NDCs[msg.sender].N2O               = HFCs;
        NDCs[msg.sender].N2O               = PFCs;
        NDCs[msg.sender].N2O               = SF6;
        NDCs[msg.sender].N2O               = NF3;
        NDCs[msg.sender].AltEnergy         = AltEnergy;
        NDCs[msg.sender].Mobilization      = Mobilization;
        NDCs[msg.sender].ContribGreenFund  = ContribGreenFund;
        NDCs[msg.sender].BilateralLoan     = BilateralLoan;
        NDCs[msg.sender].timeTarget        = timeTarget;
        NDCs[msg.sender].index             = countryIndex.push(msg.sender)-1;
        return countryIndex.length-1;
    }

    /**
       Function to update a NDC
    */
    function updateNDC(
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
        uint timeTarget)
    public
    returns(bool success)
    {
        if(!isCounrty(msg.sender)) revert();
        require(NDCs[msg.sender].country_name == country_name);
        NDCs[msg.sender].CO2               = CO2;
        NDCs[msg.sender].CH4               = CH4;
        NDCs[msg.sender].N2O               = N2O;
        NDCs[msg.sender].HFCs              = HFCs;
        NDCs[msg.sender].PFCs              = PFCs;
        NDCs[msg.sender].SF6               = SF6;
        NDCs[msg.sender].NF3               = NF3;
        NDCs[msg.sender].AltEnergy         = AltEnergy;
        NDCs[msg.sender].Mobilization      = Mobilization;
        NDCs[msg.sender].ContribGreenFund  = ContribGreenFund;
        NDCs[msg.sender].BilateralLoan     = BilateralLoan;
        NDCs[msg.sender].timeTarget        = timeTarget;
        return true;
    }

    /**
    Functions returns country NDC
    */
    function getGHG(address countryAddr)
    public
    view
    returns(
        bytes32 country_name,
        int CO2,
        int CH4,
        int N2O,
        uint timeTarget)
    {
        if(!isCounrty(countryAddr)) revert();
        return(
        NDCs[countryAddr].country_name,
        NDCs[countryAddr].CO2,
        NDCs[countryAddr].CH4,
        NDCs[countryAddr].N2O,
        NDCs[countryAddr].timeTarget
        );
    }
    function getNdcCO2(address countryAddr)
    public
    view
    returns(
        int CO2,
        uint timeTarget)
    {
        if(!isCounrty(countryAddr)) revert();
        return(
        NDCs[countryAddr].CO2,
        NDCs[countryAddr].timeTarget
        );
    }
}
