package subscriber

import (
	"context"
	"errors"
	"testing"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/pagination"
)

func TestValidateAccountNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// ===== Позитивные тесты =====

		// КЭ-2: допустимая длина, КЭ-4: допустимый алфавит,
		// КЭ-9: все типы, КЭ-13: нет подряд дефисов,
		// КЭ-15: начинается с буквы, КЭ-18: заканчивается цифрой
		{
			name:    "valid_min_length_4",
			input:   "А-Б1",
			wantErr: nil,
		},
		{
			name:    "valid_max_length_20",
			input:   "АБВГ-ДЕЁЖ-ЗИЙК-ЛМНО0",
			wantErr: nil,
		},
		{
			name:    "valid_length_5",
			input:   "АБ-В1",
			wantErr: nil,
		},
		{
			name:    "valid_length_19",
			input:   "АБВГ-ДЕЁЖ-ЗИЙК-ЛМН0",
			wantErr: nil,
		},
		{
			name:    "valid_typical_value",
			input:   "АБВ-ГДЕ-ЖЗИ-1",
			wantErr: nil,
		},
		{
			name:    "valid_multiple_non_consecutive_hyphens",
			input:   "А-Б-В-1",
			wantErr: nil,
		},
		{
			name:    "valid_with_letter_Ё",
			input:   "А-Ё1",
			wantErr: nil,
		},
		{
			name:    "valid_single_hyphen",
			input:   "АБ-1",
			wantErr: nil,
		},
		{
			name:    "valid_end_of_alphabet_letters",
			input:   "ЯЮ-ЭЩ-9",
			wantErr: nil,
		},

		// ===== Негативные: длина (КЭ-1, КЭ-3) =====

		{
			name:    "too_short_empty_string",
			input:   "",
			wantErr: ErrAccountNumberTooShort,
		},
		{
			name:    "too_short_1_char",
			input:   "А",
			wantErr: ErrAccountNumberTooShort,
		},
		{
			name:    "too_short_3_chars",
			input:   "А-1",
			wantErr: ErrAccountNumberTooShort,
		},
		{
			name:    "too_long_21_chars",
			input:   "АБВГД-ЕЖЗИК-ЛМНОП-РС0",
			wantErr: ErrAccountNumberTooLong,
		},

		// ===== Негативные: недопустимые символы (КЭ-5, КЭ-6, КЭ-7, КЭ-8) =====

		{
			name:    "lowercase_russian_letter",
			input:   "АБ-в1",
			wantErr: ErrAccountNumberInvalidChar,
		},
		{
			name:    "latin_letter_in_middle",
			input:   "АW-Б1",
			wantErr: ErrAccountNumberInvalidChar,
		},
		{
			name:    "special_char_at",
			input:   "А@Б-1",
			wantErr: ErrAccountNumberInvalidChar,
		},
		{
			name:    "special_char_underscore",
			input:   "А_Б-1",
			wantErr: ErrAccountNumberInvalidChar,
		},
		{
			name:    "special_char_slash",
			input:   "А/Б-1",
			wantErr: ErrAccountNumberInvalidChar,
		},
		{
			name:    "space_in_middle",
			input:   "А Б-1",
			wantErr: ErrAccountNumberInvalidChar,
		},

		// ===== Негативные: нет дефисов (КЭ-12) =====

		{
			name:    "no_hyphens",
			input:   "АБВГ1",
			wantErr: ErrAccountNumberNoHyphen,
		},

		// ===== Негативные: подряд идущие дефисы (КЭ-14) =====

		{
			name:    "two_consecutive_hyphens",
			input:   "А--Б1",
			wantErr: ErrAccountNumberConsecutiveHyphens,
		},
		{
			name:    "three_consecutive_hyphens",
			input:   "А---Б1",
			wantErr: ErrAccountNumberConsecutiveHyphens,
		},

		// ===== Негативные: первый символ (КЭ-16, КЭ-17) =====

		{
			name:    "starts_with_digit",
			input:   "1А-Б2",
			wantErr: ErrAccountNumberStartNotLetter,
		},
		{
			name:    "starts_with_hyphen",
			input:   "-АБВ1",
			wantErr: ErrAccountNumberStartNotLetter,
		},
		{
			name:    "starts_with_latin_letter",
			input:   "WА-Б1",
			wantErr: ErrAccountNumberStartNotLetter,
		},

		// ===== Негативные: последний символ (КЭ-19, КЭ-20) =====

		{
			name:    "ends_with_letter",
			input:   "А-1Б",
			wantErr: ErrAccountNumberEndNotDigit,
		},
		{
			name:    "ends_with_hyphen",
			input:   "АБ-1-",
			wantErr: ErrAccountNumberEndNotDigit,
		},

		// ===== Комбинированные нарушения =====

		{
			name:    "multiple_violations_digit_start_and_consecutive_hyphens",
			input:   "1--Б2",
			wantErr: ErrAccountNumberStartNotLetter,
		},
		{
			name:    "multiple_violations_no_hyphens_and_ends_with_letter",
			input:   "АБВГД",
			wantErr: ErrAccountNumberEndNotDigit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAccountNumber(tt.input)

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("ValidateAccountNumber(%q) = %v, want nil", tt.input, err)
				}
				return
			}

			if err == nil {
				t.Errorf("ValidateAccountNumber(%q) = nil, want %v", tt.input, tt.wantErr)
				return
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateAccountNumber(%q) = %v, want %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

// mockRepository is a test double for the Repository interface.
type mockRepository struct {
	addSubscriberFunc           func(ctx context.Context, request AddSubscriberRequest) (Subscriber, error)
	getSubscriberByIDFunc       func(ctx context.Context, id int) (Subscriber, error)
	getSubscriberExtendedIDFunc func(ctx context.Context, id int) (ExtendedSubscriber, error)
	getAllSubscribersFunc       func(ctx context.Context, page pagination.Pagination) ([]Subscriber, error)
}

func (m *mockRepository) AddSubscriber(ctx context.Context, request AddSubscriberRequest) (Subscriber, error) {
	if m.addSubscriberFunc != nil {
		return m.addSubscriberFunc(ctx, request)
	}
	return Subscriber{ID: 1, AccountNumber: request.AccountNumber}, nil
}

func (m *mockRepository) GetSubscriberByID(ctx context.Context, id int) (Subscriber, error) {
	if m.getSubscriberByIDFunc != nil {
		return m.getSubscriberByIDFunc(ctx, id)
	}
	return Subscriber{}, nil
}

func (m *mockRepository) GetSubscriberExtendedByID(ctx context.Context, id int) (ExtendedSubscriber, error) {
	if m.getSubscriberExtendedIDFunc != nil {
		return m.getSubscriberExtendedIDFunc(ctx, id)
	}
	return ExtendedSubscriber{}, nil
}

func (m *mockRepository) GetAllSubscribers(ctx context.Context, page pagination.Pagination) ([]Subscriber, error) {
	if m.getAllSubscribersFunc != nil {
		return m.getAllSubscribersFunc(ctx, page)
	}
	return nil, nil
}

func TestAddSubscriber_ValidAccountNumber(t *testing.T) {
	repo := &mockRepository{}
	svc := NewService(repo)

	ctx := goctx.Wrap(context.Background())

	result, err := svc.AddSubscriber(ctx, AddSubscriberRequest{
		AccountNumber: "АБВ-ГДЕ-1",
	})

	if err != nil {
		t.Fatalf("AddSubscriber() unexpected error: %v", err)
	}
	if result.AccountNumber != "АБВ-ГДЕ-1" {
		t.Errorf("AddSubscriber() AccountNumber = %q, want %q", result.AccountNumber, "АБВ-ГДЕ-1")
	}
}

func TestAddSubscriber_InvalidAccountNumber(t *testing.T) {
	repoCalled := false
	repo := &mockRepository{
		addSubscriberFunc: func(ctx context.Context, request AddSubscriberRequest) (Subscriber, error) {
			repoCalled = true
			return Subscriber{}, nil
		},
	}
	svc := NewService(repo)

	ctx := goctx.Wrap(context.Background())

	_, err := svc.AddSubscriber(ctx, AddSubscriberRequest{
		AccountNumber: "invalid",
	})

	if err == nil {
		t.Fatal("AddSubscriber() expected error for invalid account number, got nil")
	}

	if repoCalled {
		t.Error("AddSubscriber() should not call repository when validation fails")
	}
}

func TestAddSubscriber_ValidationErrorWrapping(t *testing.T) {
	repo := &mockRepository{}
	svc := NewService(repo)

	ctx := goctx.Wrap(context.Background())

	_, err := svc.AddSubscriber(ctx, AddSubscriberRequest{
		AccountNumber: "АБВГ1",
	})

	if err == nil {
		t.Fatal("AddSubscriber() expected error, got nil")
	}

	if !errors.Is(err, ErrAccountNumberNoHyphen) {
		t.Errorf("AddSubscriber() error = %v, want wrapped %v", err, ErrAccountNumberNoHyphen)
	}
}
