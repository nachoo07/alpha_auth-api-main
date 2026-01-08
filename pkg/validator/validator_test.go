package validator

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	err := RegisterValidation()
	assert.Nil(t, err)

	type entity struct {
		Required        string `validate:"required"`
		Number          string `validate:"numeric"`
		ValidSiteFormat string `validate:"validSiteFormat"`
		ValidStepType   string `validate:"validStepTypeFormat"`
		ValidTime       string `validate:"validTimeFormat"`
		ValidChannel    string `validate:"validChannelFormat"`
		ValidModal      string `validate:"validModalFormat"`
	}

	tests := []struct {
		name   string
		entity interface{}
		err    error
	}{
		{
			name: "given an valid entity, should not return error",
			entity: entity{
				Required:        "value1",
				Number:          "1200",
				ValidSiteFormat: "MLM",
				ValidStepType:   "first_mile",
				ValidTime:       "12:00",
				ValidChannel:    "commercial",
				ValidModal:      "air",
			},
			err: nil,
		},
		{
			name: "given an invalid entity, should return error",
			entity: entity{
				Required:        "",
				Number:          "abc",
				ValidSiteFormat: "",
				ValidStepType:   "step",
				ValidTime:       "25:00",
				ValidChannel:    "channel",
				ValidModal:      "modal",
			},
			err: fmt.Errorf(
				"Key: 'entity.Required' Error:Field validation for 'Required' failed on the 'required' tag" +
					"\nKey: 'entity.Number' Error:Field validation for 'Number' failed on the 'numeric' tag" +
					"\nKey: 'entity.ValidSiteFormat' Error:Field validation for 'ValidSiteFormat' failed on the 'validSiteFormat' tag" +
					"\nKey: 'entity.ValidStepType' Error:Field validation for 'ValidStepType' failed on the 'validStepTypeFormat' tag" +
					"\nKey: 'entity.ValidTime' Error:Field validation for 'ValidTime' failed on the 'validTimeFormat' tag" +
					"\nKey: 'entity.ValidChannel' Error:Field validation for 'ValidChannel' failed on the 'validChannelFormat' tag" +
					"\nKey: 'entity.ValidModal' Error:Field validation for 'ValidModal' failed on the 'validModalFormat' tag",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.entity)

			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestIsValidSiteFormat(t *testing.T) {
	tests := []struct {
		name      string
		sitesID   []string
		wantValid bool
	}{
		{
			name:      "given all valid sites format, all must be valid",
			sitesID:   []string{"MLA", "MLB", "MLC", "MLD", "MLX", "MLZ", "ABC"},
			wantValid: true,
		},
		{
			name:      "given some invalid sites, all must be invalid",
			sitesID:   []string{"mla", "mlb", "mlm", "mlc", "mld", "mle", "mlz", "M10", "MLAB", "MLB1"},
			wantValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, siteID := range tt.sitesID {
				got := IsValidSiteFormat(siteID)
				assert.Equal(t, tt.wantValid, got)
			}
		})
	}
}

func TestRegisterValidationMap(t *testing.T) {
	validFunc := func(validator.FieldLevel) bool { return true }
	tests := []struct {
		name       string
		tag        string
		validation func(validator.FieldLevel) bool
		wantErr    bool
	}{
		{
			name:       "given a valid entry, should not return error",
			tag:        "isValid",
			validation: validFunc,
			wantErr:    false,
		},
		{
			name:       "given an invalid entry, should return error",
			tag:        "",
			validation: validFunc,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vMap := map[string]func(validator.FieldLevel) bool{
				tt.tag: tt.validation,
			}
			err := registerValidationMap(vMap)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestIsValidTimeFormat(t *testing.T) {
	type args struct {
		time string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should be a valid identifier",
			args: args{time: ""},
			want: true,
		},
		{
			name: "should be a valid for 24hr format last minute",
			args: args{time: "23:59"},
			want: true,
		},
		{
			name: "should be a valid for 24hr format first minute",
			args: args{time: "00:00"},
			want: true,
		},
		{
			name: "should not be a valid for 24hr format above the limit",
			args: args{time: "24:00"},
			want: false,
		},
		{
			name: "should be not valid if missing digits",
			args: args{time: "2:0"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidTimeFormat(tt.args.time); got != tt.want {
				t.Errorf("IsValidTimeFormat(%s) = %v, want %v", tt.args.time, got, tt.want)
			}
		})
	}
}

func TestIsValidStepTypeFormat(t *testing.T) {
	type args struct {
		val string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should not be valid step type",
			args: args{val: "line_haul"},
			want: false,
		},
		{
			name: "should be valid for middle_mile",
			args: args{val: "middle_mile"},
			want: true,
		},
		{
			name: "should be valid for first_mile",
			args: args{val: "first_mile"},
			want: true,
		},
		{
			name: "should be valid for last_mile",
			args: args{val: "last_mile"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidStepTypeFormat(tt.args.val); got != tt.want {
				t.Errorf("IsValidStepTypeFormat(%s) = %v, want %v", tt.args.val, got, tt.want)
			}
		})
	}
}

func TestIsValidChannelFormat(t *testing.T) {
	type args struct {
		val string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should not be valid channel",
			args: args{val: "channel"},
			want: false,
		},
		{
			name: "should be valid for commercial",
			args: args{val: "commercial"},
			want: true,
		},
		{
			name: "should be valid for logistics",
			args: args{val: "logistics"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidChannelFormat(tt.args.val); got != tt.want {
				t.Errorf("IsValidChannelFormat(%s) = %v, want %v", tt.args.val, got, tt.want)
			}
		})
	}
}

func TestIsValidModalFormat(t *testing.T) {
	type args struct {
		val string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should not be valid modal",
			args: args{val: "modal"},
			want: false,
		},
		{
			name: "should be valid for air",
			args: args{val: "air"},
			want: true,
		},
		{
			name: "should be valid for ground",
			args: args{val: "ground"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidModalFormat(tt.args.val); got != tt.want {
				t.Errorf("IsValidModalFormat(%s) = %v, want %v", tt.args.val, got, tt.want)
			}
		})
	}
}

func TestIsValidDateFormat(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "should be valid date",
			value: "2020-12-31",
			want:  true,
		},
		{
			name:  "should not be valid date",
			value: "2020-13-31",
			want:  false,
		},
		{
			name:  "should be valid empty string",
			value: "",
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidDateFormat(tt.value); got != tt.want {
				t.Errorf("IsValidDateFormat(%s) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestIsDateFromBeforeOrEqualDateTo(t *testing.T) {
	dateFrom := "2020-12-31"
	dateTo := "2021-01-01"
	dateFromTime, _ := time.Parse(time.DateOnly, dateFrom)
	dateToTime, _ := time.Parse(time.DateOnly, dateTo)
	tests := []struct {
		name     string
		dateFrom time.Time
		dateTo   time.Time
		want     bool
	}{
		{
			name:     "should be equal dates",
			dateFrom: dateFromTime,
			dateTo:   dateFromTime,
			want:     true,
		},

		{
			name:     "should be valid interval",
			dateFrom: dateFromTime,
			dateTo:   dateToTime,
			want:     true,
		},
		{
			name:     "should not be valid interval",
			dateFrom: dateToTime,
			dateTo:   dateFromTime,
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDateFromBeforeOrEqualDateTo(tt.dateFrom, tt.dateTo); got != tt.want {
				t.Errorf("IsDateFromBeforeOrEqualDateTo(%s, %s) = %v, want %v", tt.dateFrom, tt.dateTo, got, tt.want)
			}
		})
	}
}
