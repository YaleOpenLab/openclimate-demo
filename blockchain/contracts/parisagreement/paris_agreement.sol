pragma solidity ^0.5.10;

contract ParisAgreementHighLevel{
    // At the high level we can set goals for every country and receive accumulated reports to keep track/actions

    /**
    Each country will be assigned is public ETH address and a private key kept secret on the responsible state agency
    We can rethink this process later.
    */
    address[] public countries;

    /**
    Global goals set by Conference of the Parties (COP). Related Voting contract will be separately added later.
    */

    struct global_stocktake {
        address[] votingPatries;
        int CO2; // (required) in metric tonnes
        int CH4; // (required) in metric tonnes
        int N2O; // (required) in metric tonnes
        int AltEnergy; // alternative/renewable energy usage in MWh
        uint timeGoal;
    }

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
        string country;
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
    Probably we would need a intermediate total GHG variables to track who we do on the global level. Updated by information incoming from databases every year or month.
    */
    struct global_reduction_track {
        int  earthCO2;
        int  earthCH4;
        int  earthN2O;
        int  earthAltEnergy;
        uint  lastReportTime;
    }

    /**
    Creating mapping to store countries goals. This allows to save/delete/update pledges for each country.
    */
    mapping(address => country_NDC) private NDCs;
    address[] private countryIndex;


    /**
    Util function to emit variable after inserting/updating/deleting
    */
    event LogNewNDC   (address indexed country_address, uint index, string country_name, int CO2, int CH4, int N2O, int HFCs, int PFCs, int SF6, int NF3, int AltEnergy, uint timeTarget);
    event LogUpdateNDC(address indexed country_address, uint index, string country_name, int CO2, int CH4, int N2O, int HFCs, int PFCs, int SF6, int NF3, int AltEnergy, uint timeTarget);
    event LogDeleteNDC(address indexed country_address, uint index);

    /**
    Util function toto check if the public key of the country in the list
    */
    function isCounrty(address counrty)
    public
    view
    returns(bool isIndeed)
    {
        if(countryIndex.length == 0) return false;
        return (countryIndex[Pledges[counrty].index] == counrty);
    }

    /**
       Function to insert a new NDC from a country address
    */
    function insertNDC(
        string memory country,
        int CO2,
        int CH4,
        int N2O,
        int HFCs,
        int PFCs,
        int SF6,
        int NF3,
        int AltEnergy,
        uint timeTarget)
    public
    returns(uint index)
    {
        if(isCounrty(msg.sender)) revert();
        NDCs[msg.sender].country           = country;
        NDCs[msg.sender].CO2               = CO2;
        NDCs[msg.sender].CH4               = CH4;
        NDCs[msg.sender].N2O               = N2O;
        NDCs[msg.sender].N2O               = HFCs;
        NDCs[msg.sender].N2O               = PFCs;
        NDCs[msg.sender].N2O               = SF6;
        NDCs[msg.sender].N2O               = NF3;
        NDCs[msg.sender].AltEnergy         = AltEnergy;
        NDCs[msg.sender].timeTarget        = timeTarget;
        NDCs[msg.sender].index             = countryIndex.push(msg.sender)-1;
        emit LogNewNDC(
            msg.sender,
            Pledges[msg.sender].index,
            country,
            CO2,
            CH4,
            N2O,
            HFCs,
            PFCs,
            SF6,
            NF3,
            AltEnergy,
                timeTarget);
        return countryIndex.length-1;
    }

    /**
       Function to update a NDC
    */
    function updateNDC(address country_address, uint index, string country_name, int CO2, int CH4, int N2O, int HFCs, int PFCs, int SF6, int NF3, int AltEnergy, uint timeTarget)
    public
    returns(bool success)
    {
        if(!isCounrty(country_address)) revert();
        require(NDCs[msg.sender].country == country_name);
            NDCs[msg.sender].CO2               = CO2;
            NDCs[msg.sender].CH4               = CH4;
            NDCs[msg.sender].N2O               = N2O;
            NDCs[msg.sender].HFCs               = HFCs;
            NDCs[msg.sender].PFCs               = PFCs;
            NDCs[msg.sender].SF6               = SF6;
            NDCs[msg.sender].NF3               = NF3;
            NDCs[msg.sender].AltEnergy         = AltEnergy;
            NDCs[msg.sender].timeTarget        = timeTarget;

            emit LogUpdateNDC(
                country_address,
                NDCs[msg.sender].index,
                NDCs[msg.sender].CO2,
                NDCs[msg.sender].CH4,
                NDCs[msg.sender].N2O,
                NDCs[msg.sender].HFCs,
                NDCs[msg.sender].PFCs,
                NDCs[msg.sender].SF6,
                NDCs[msg.sender].NF3,
                NDCs[msg.sender].AltEnergy,
                NDCs[msg.sender].timeTarget);
            return true;
    }

    /**
     Util function to check time left to a time goal set up by COP (indicatively every 5 years)
     */
    function checkDeadline () public returns (uint time, bool overdue){
        uint timeLeft = global_stocktake.timeGoal - now;
        if (timeLeft >= 0){
            return (timeLeft, false);
        }
        return (timeLeft, true);
    }





}
