package types

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	govtypes "github.com/okex/okexchain/x/gov/types"
	"github.com/stretchr/testify/suite"
)

const (
	expectedTitle                                           = "default title"
	expectedDescription                                     = "default description"
	expectedManageContractDeploymentWhitelistProposalString = `ManageContractDeploymentWhitelistProposal:
 Title:					default title
 Description:        	default description
 Type:                	ManageContractDeploymentWhitelist
 IsAdded:				true
 DistributorAddrs:
						okexchain1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqupa6dx
						okexchain1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqpphf0s5`
)

type ProposalTestSuite struct {
	suite.Suite
	strBuilder strings.Builder
	addrs      []sdk.AccAddress
}

func TestProposalTestSuite(t *testing.T) {
	proposalTestSuite := ProposalTestSuite{
		addrs: []sdk.AccAddress{
			ethcmn.BytesToAddress([]byte{0x0}).Bytes(),
			ethcmn.BytesToAddress([]byte{0x1}).Bytes(),
		},
	}
	suite.Run(t, &proposalTestSuite)
}

func (suite *ProposalTestSuite) TestProposal_ManageContractDeploymentWhitelistProposal() {
	proposal := NewManageContractDeploymentWhitelistProposal(
		expectedTitle,
		expectedDescription,
		suite.addrs,
		true,
	)

	suite.Require().Equal(expectedTitle, proposal.GetTitle())
	suite.Require().Equal(expectedDescription, proposal.GetDescription())
	suite.Require().Equal(RouterKey, proposal.ProposalRoute())
	suite.Require().Equal(proposalTypeManageContractDeploymentWhitelist, proposal.ProposalType())
	suite.Require().Equal(expectedManageContractDeploymentWhitelistProposalString, proposal.String())

	testCases := []struct {
		msg           string
		prepare       func()
		expectedError bool
	}{
		{
			"pass",
			func() {},
			false,
		},
		{
			"empty title",
			func() {
				proposal.Title = ""
			},
			true,
		},
		{
			"overlong title",
			func() {
				for i := 0; i < govtypes.MaxTitleLength+1; i++ {
					suite.strBuilder.WriteByte('a')
				}
				proposal.Title = suite.strBuilder.String()
			},
			true,
		},
		{
			"empty description",
			func() {
				proposal.Description = ""
				proposal.Title = expectedTitle
			},
			true,
		},
		{
			"overlong description",
			func() {
				suite.strBuilder.Reset()
				for i := 0; i < govtypes.MaxDescriptionLength+1; i++ {
					suite.strBuilder.WriteByte('a')
				}
				proposal.Description = suite.strBuilder.String()
			},
			true,
		},
		{
			"duplicated distributor addresses",
			func() {
				// add a duplicated address into DistributorAddrs
				proposal.DistributorAddrs = append(proposal.DistributorAddrs, proposal.DistributorAddrs[0])
				proposal.Description = expectedDescription
			},
			true,
		},
		{
			"empty distributor addresses",
			func() {
				proposal.DistributorAddrs = nil
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			tc.prepare()

			err := proposal.ValidateBasic()

			if tc.expectedError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
